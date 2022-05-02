package route

import (
	"fmt"
	"github.com/M-Cosmosss/mclib/internal/context"
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/M-Cosmosss/mclib/internal/utils"
	"github.com/pkg/errors"
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
	err := db.Books.Create(ctx.Request().Context(), db.CreateBookOptions{
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

func (b *Books) Delete(ctx *context.Context) error {
	id := ctx.QueryInt("id")
	_, err := db.Books.GetByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.Error(http.StatusNotFound, db.ErrBookNotExist.Error())
	}

	err = db.Books.Delete(ctx.Request().Context(), id)
	if err != nil {
		log.Error("Delete book failed: %v", err)
		return ctx.Error(http.StatusInternalServerError, err.Error())
	}
	log.Info("Delete book: %s", id)
	return ctx.Success("删除成功")
}

//
//type ListBookOption struct {
//	ID  uint `form:"id" binding:"Required"`
//	Num int  `form:"num" binding:"Required"`
//}

func (b *Books) List(ctx *context.Context) error {
	off, num := ctx.QueryInt("offset"), ctx.QueryInt("num")
	if num <= 0 {
		return ctx.Error(http.StatusBadRequest, "错误参数")
	}
	fmt.Printf("id-%d num-%d\n", off, num)
	books, count, err := db.Books.List(ctx.Request().Context(), db.ListBookOption{
		Off: off,
		Num: num,
	})
	if err != nil {
		return ctx.Error(http.StatusInternalServerError, err.Error())
	}

	return ctx.Success(map[string]interface{}{
		"count": count,
		"books": books,
	})
}

func (b *Books) Get(ctx *context.Context) error {
	id := uint(ctx.QueryInt("id"))
	isbn := ctx.QueryInt("isbn")
	name := ctx.Query("name")
	author := ctx.Query("author")

	book, err := db.Books.Get(ctx.Request().Context(), db.GetBookOptions{
		ID:     id,
		ISBN:   isbn,
		Name:   name,
		Author: author,
	})
	if err != nil {
		return ctx.Error(http.StatusNotFound, "书籍不存在")
		//}
		//log.Error("Failed to get book: %v", err)
		//return ctx.Error(http.StatusInternalServerError, "内部错误")
	}

	return ctx.Success(book)
}

type GetBookOption struct {
	ID int `json:"id"`
}

type GetBookDetailReply struct {
	Name       string
	Author     string
	Logo       string
	ISBN       int
	IsBorrowed bool
	Logs       []db.BorrowLog
}

func (b *Books) GetDetailByID(ctx *context.Context) error {
	id := ctx.ParamInt("id")
	book, err := db.Books.GetByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.Error(http.StatusNotFound, err.Error())
	}
	logs, err := db.BorrowLogs.GetByBookID(ctx.Request().Context(), db.GetLogByBookIDOption{
		BookID: id,
		Offset: 0,
		Limit:  5,
	})
	reply := &GetBookDetailReply{
		Name:       book.Name,
		Author:     book.Author,
		Logo:       book.Logo,
		ISBN:       book.ISBN,
		IsBorrowed: book.IsBorrowed,
		Logs:       logs,
	}
	return ctx.Success(reply)
}

type BorrowBookOption struct {
	BookID int `json:"book_id"`
}

func (b *Books) Borrow(ctx *context.Context, user *db.User, option BorrowBookOption) error {
	if len(user.BooksID) >= user.BooksLimit {
		return ctx.Error(http.StatusForbidden, "用户借书量超过限制")
	}
	if err := db.Books.Borrow(ctx.Request().Context(), option.BookID, int(user.ID)); err != nil {
		return ctx.Error(http.StatusForbidden, err.Error())
	}
	_ = db.BorrowLogs.Create(ctx.Request().Context(), db.CreateLogOption{
		User:      int(user.ID),
		Book:      option.BookID,
		Operation: db.BORROW,
	})

	if err := db.Users.Borrow(ctx.Request().Context(), *user, option.BookID); err != nil {
		return ctx.Error(http.StatusInternalServerError, fmt.Sprintf("数据库错误： %v", err))
	}

	return ctx.Success("借书成功")
}

func (b *Books) Return(ctx *context.Context) error {
	id := ctx.QueryInt("id")
	if book, err := db.Books.GetByID(ctx.Request().Context(), int(id)); err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error())
	} else {
		if !book.IsBorrowed {
			return ctx.Error(http.StatusBadRequest, "该书未处于被借状态")
		}

		if user, err := db.Users.GetByID(ctx.Request().Context(), uint(book.LastBorrowedUserID)); err != nil {
			return ctx.Error(http.StatusInternalServerError, fmt.Sprintf("数据库错误，请联系管理员: %v", err))
		} else {
			nums, err := utils.CleanNumInSlice(user.BooksID, int32(id))
			if err != nil {
				return ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "CleanNumInSlice").Error())
			}
			if err := db.Users.Return(ctx.Request().Context(), *user, nums); err != nil {
				return ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "user return").Error())
			}
			if err := db.BorrowLogs.Create(ctx.Request().Context(), db.CreateLogOption{
				User:      int(user.ID),
				Book:      int(id),
				Operation: db.RETURN,
			}); err != nil {
				return ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "borrow log").Error())
			}
			book.IsBorrowed = false
			if err := db.Books.UpdateByID(ctx.Request().Context(), *book); err != nil {
				return ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "update book").Error())
			}
			return ctx.Success("还书成功")
		}
	}
}

func (b *Books) Logs(ctx *context.Context) error {
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")
	borrowLogs, count, err := db.BorrowLogs.List(ctx.Request().Context(), db.ListBorrowLogsOptions{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Error("Failed to list borrow logs: %v", err)
		return ctx.Error(http.StatusInternalServerError, "内部错误")
	}
	return ctx.Success(map[string]interface{}{
		"borrow_logs": borrowLogs,
		"count":       count,
	})
}
