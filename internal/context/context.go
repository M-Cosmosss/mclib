package context

import (
	"encoding/json"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/flamego/flamego"
	"log"
	"net/http"

	"github.com/flamego/session"
)

type Context struct {
	flamego.Context
	IsUser    bool
	IsManager bool
}

func (c *Context) Success(data ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	err := json.NewEncoder(c.ResponseWriter()).Encode(
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
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(errorCode)

	err := json.NewEncoder(c.ResponseWriter()).Encode(
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

func Contexter() flamego.Handler {
	return func(ctx flamego.Context, sess session.Session) {
		c := &Context{ctx, false, false}
		userID, ok := sess.Get(UserIDSessionKey).(uint)
		if ok {
			u, err := db.Users.GetByID(ctx.Request().Context(), userID)
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
