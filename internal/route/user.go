package route

import (
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/flamego/session"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
	user, err := db.Users.Create(ctx.Request().Context(), opts.CreateUserOption)
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
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *Users) Login(ctx *context.Context, sess session.Session, opts UserLoginOption) error {
	user, err := db.Users.GetByName(ctx.Request().Context(), opts.Name)
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error())
	}
	if user.ValidatePassword(opts.Password) {
		sess.Set(context.UserIDSessionKey, user.ID)
		return ctx.Success(map[string]interface{}{
			"isAdmin": user.IsManager,
			"token":   sess.ID(),
		})
	}
	log.Info("Login failed")
	return ctx.Error(http.StatusUnauthorized, "账号或密码错误")
}

func (u *Users) Logout(ctx *context.Context, sess session.Session) error {
	sess.Flush()
	return ctx.Success("登出成功")
}

type UserReply struct {
	Name       string
	BooksLimit int
	BooksID    []int32
}

func (u *Users) GetSelf(ctx *context.Context, user *db.User) error {
	reply := &UserReply{
		Name:       user.Name,
		BooksLimit: user.BooksLimit,
		BooksID:    user.BooksID,
	}
	return ctx.Success(reply)
}

type ManagerUserReply struct {
	gorm.Model

	Name       string `gorm:"unique"`
	BooksLimit int
	IsManager  bool
	BooksID    pq.Int32Array `gorm:"type:integer[]"`
}

func (u *Users) GetByID(ctx *context.Context) error {
	id := ctx.ParamInt("id")
	var user *db.User
	var err error
	if user, err = db.Users.GetByID(ctx.Request().Context(), uint(id)); err != nil {
		return ctx.Error(http.StatusNotFound, db.ErrUserNotExists.Error())
	}
	return ctx.Success(&ManagerUserReply{
		Model:      user.Model,
		Name:       user.Name,
		BooksLimit: user.BooksLimit,
		IsManager:  user.IsManager,
		BooksID:    user.BooksID,
	})
}

func (u *Users) DeleteByID(ctx *context.Context) error {
	id := ctx.ParamInt("id")
	err := db.Users.DeleteByID(ctx.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, db.ErrBookNotExist) {
			return ctx.Error(http.StatusNotFound, "用户不存在")
		}
		log.Error("Failed to delete user by ID: %v", err)
		return ctx.Error(http.StatusInternalServerError, "系统内部错误")
	}
	return ctx.Success("删除用户成功")
}
