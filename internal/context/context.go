package context

import (
	"encoding/json"
	"github.com/M-Cosmosss/mclib/internal/db"
	"log"
	"net/http"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type Context struct {
	*macaron.Context
	IsUser    bool
	IsManager bool
}

func (c *Context) Success(data ...interface{}) error {
	c.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Resp.WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	err := json.NewEncoder(c.Resp).Encode(
		map[string]interface{}{
			"error": 0,
			"data":  d,
		},
	)
	if err != nil {
		log.Printf("Failed to encode: %v", err)
	}
	return err
}

func (c *Context) Error(errorCode int, msg string) error {
	c.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Resp.WriteHeader(errorCode)

	err := json.NewEncoder(c.Resp).Encode(
		map[string]interface{}{
			"error": errorCode,
			"msg":   msg,
		},
	)
	if err != nil {
		log.Printf("Failed to encode: %v", err)
	}
	return nil
}

func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store) {
		c := &Context{ctx, false, false}
		userID, ok := sess.Get(UserIDSessionKey).(uint)
		if ok {
			u, err := db.Users.GetByID(ctx.Req.Context(), userID)
			if err == nil {
				c.IsUser = true
			}
			if err == nil && u.IsManager {
				c.IsManager = true
			}
			c.Map(u)
		}

		c.Map(c)
	}
}
