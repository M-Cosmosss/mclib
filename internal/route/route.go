package route

import (
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/session/postgres"
	"net/http"
)

func Init() *flamego.Flame {
	m := flamego.Classic()

	var sessionStorage interface{}
	initer := postgres.Initer()
	sessionStorage = postgres.Config{
		DSN: db.GetPGDSN(),
	}

	sessioner := session.Sessioner(session.Options{
		Initer:      initer,
		Config:      sessionStorage,
		ReadIDFunc:  func(r *http.Request) string { return r.Header.Get("Authorization") },
		WriteIDFunc: func(w http.ResponseWriter, r *http.Request, sid string, created bool) {},
	})

	m.Use(sessioner)
	m.Use(context.Contexter())

	m.Group("/user", func() {
		m.Post("/register", binding.JSON(CreateUserOption{}), UserHandler.Create)
		m.Post("/login", binding.JSON(UserLoginOption{}), UserHandler.Login)

		m.Group("", func() {
			m.Get("/self", UserHandler.GetSelf)
			m.Post("/logout", UserHandler.Logout)
		}, userAuth)

		m.Group("", func() {
			m.Group("/id/{id}", func() {
				m.Get("", UserHandler.GetByID)
				m.Delete("", UserHandler.DeleteByID)
			})
			m.Get("/logs", BookHandler.Logs)
		}, managerAuth)
	})
	//m.Use(userAuth)
	m.Group("/book", func() {
		m.Get("/list", BookHandler.List)
		m.Get("", BookHandler.Get)

		m.Post("/borrow", userAuth, binding.JSON(BorrowBookOption{}), BookHandler.Borrow)
		m.Delete("/borrow", managerAuth, BookHandler.Return)

		m.Group("", func() {
			m.Post("", binding.JSON(CreateBookOption{}), BookHandler.Create)
			m.Delete("", BookHandler.Delete)
		}, managerAuth)
	})

	m.Get("/test", userAuth, func(ctx *context.Context, user *db.User) error {
		return ctx.Success(map[string]interface{}{
			"msg":       "success",
			"isManager": user.IsManager,
		})
	})
	return m
}
