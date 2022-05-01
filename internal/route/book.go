package route

import (
	"fmt"
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"net/http"
	log "unknwon.dev/clog/v2"
)

var BookHandler = &Books{}

type Books struct {
}

type CreateBookOption struct {
	Name   string `json:"name" binding:"Required"`
	Author string `json:"author" binding:"Required"`
	Logo   string `json:"logo" binding:"Required"`
	ISBN   int    `json:"isbn" binding:"Required"`
}

func (b *Books) Create(ctx *context.Context, option CreateBookOption) error {
	//fmt.Printf("%+v", option)
	err := db.Books.Create(ctx.Req.Context(), db.CreateBookOptions{
		Name:   option.Name,
		Author: option.Author,
		ISBN:   option.ISBN,
		Logo:   option.Logo,
	})
	if err != nil {
		log.Error("Create book failed: %v", err)
		return ctx.Error(http.StatusInternalServerError, err.Error())
	}
	log.Info("Create book: %s", option.Name)
	return nil
}

type ListBookOption struct {
	ID  uint `form:"id" binding:"Required"`
	Num int  `form:"num" binding:"Required"`
}

func (b *Books) List(ctx *context.Context) error {
	id, num := ctx.QueryInt("id"), ctx.QueryInt("num")
	if num <= 0 {
		return ctx.Error(http.StatusBadRequest, "错误参数")
	}
	fmt.Printf("id-%d num-%d\n", id, num)
	books, err := db.Books.List(ctx.Req.Context(), db.ListBookOption{
		ID:  uint(id),
		Num: num,
	})
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error())
	}
	return ctx.Success(books)
}

type BorrowBookOption struct {
	BookID uint
}

func (b *Books) Borrow(ctx *context.Context, user *db.User, option BorrowBookOption) error {
	if len(user.BooksID) >= user.BooksLimit {
		return ctx.Error(http.StatusForbidden, "用户借书量超过限制")
	}
	if err := db.Books.Borrow(ctx.Req.Context(), option.BookID); err != nil {
		return ctx.Error(http.StatusForbidden, err.Error())
	}

}
