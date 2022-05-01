package route

import (
	"fmt"
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/go-macaron/session"
	"net/http"
	log "unknwon.dev/clog/v2"
)

var UserHandler = &Users{}

type Users struct {
}

type CreateUserOption struct {
	db.CreateUserOption
}

func (u *Users) Create(ctx *context.Context, opts CreateUserOption) error {
	user, err := db.Users.Create(ctx.Req.Context(), opts.CreateUserOption)
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error())
	}
	return ctx.Success(struct {
		ID   uint
		Name string
	}{
		ID:   user.ID,
		Name: user.Name,
	})
}

type UserLoginOption struct {
	Name     string
	Password string
}

func (u *Users) Login(ctx *context.Context, sess session.Store, opts UserLoginOption) error {
	user, err := db.Users.GetByName(ctx.Req.Context(), opts.Name)
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error())
	}
	if user.ValidatePassword(opts.Password) {
		_ = sess.Set(context.UserIDSessionKey, user.ID)
		return ctx.Success(map[string]string{
			"isAdmin": fmt.Sprintf("%t", user.IsManager),
			"token":   sess.ID(),
		})
	}
	log.Info("Login failed")
	return ctx.Error(http.StatusUnauthorized, "账号或密码错误")
}

func (u *Users) Borrow(ctx *context.Context, sess session.Store, opts UserLoginOption) error {
	//user, err := db.Users.GetByName(ctx.Req.Context(), opts.Name)
	//if err != nil {
	//	return ctx.Error(http.StatusBadRequest, err.Error())
	//}
	//if user.ValidatePassword(opts.Password) {
	//	_ = sess.Set(context.UserIDSessionKey, user.ID)
	//	return ctx.Success(map[string]string{"msg": "Login success"})
	//}
	//log.Info("Login failed")
	return nil
}
