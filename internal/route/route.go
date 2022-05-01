package route

import (
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func Init() *macaron.Macaron {
	m := macaron.Classic()
	m.Use(session.Sessioner())
	m.Use(context.Contexter())

	m.Group("/user", func() {
		m.Post("/register", binding.Bind(CreateUserOption{}), UserHandler.Create)
		m.Post("/login", binding.Bind(UserLoginOption{}), UserHandler.Login)
	})
	m.Get("/test", userAuth, func(ctx *context.Context, user *db.User) error {
		return ctx.Success(map[string]interface{}{
			"msg":       "success",
			"isManager": user.IsManager,
		})
	})
	//m.Use(userAuth)
	m.Group("/book", func() {
		m.Get("/list", BookHandler.List)
		m.Post("/create", binding.Bind(CreateBookOption{}), BookHandler.Create)
		m.Post("/borrow", binding.Bind(BorrowBookOption{}), BookHandler.Borrow)
	})
	return m
}
