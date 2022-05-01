package db

import (
	"context"
	"github.com/pkg/errors"
	"time"
	log "unknwon.dev/clog/v2"

	"gorm.io/gorm"
)

var _ BooksStore = (*books)(nil)

var Books BooksStore

type BooksStore interface {
	Create(ctx context.Context, opts CreateBookOptions) error
	Get(ctx context.Context, opts GetBookOptions) (*Book, error)
	GetByID(ctx context.Context, id int) (*Book, error)
	GetByISBN(ctx context.Context, isbn int) ([]Book, error)
	List(ctx context.Context, option ListBookOption) ([]Book, int, error)
	Borrow(ctx context.Context, bid int, uid int) error
	Delete(ctx context.Context, id int) error
}

type books struct {
	*gorm.DB
}

type Book struct {
	gorm.Model

	Name               string
	Author             string
	Logo               string
	ISBN               int
	IsBorrowed         bool
	LastBorrowedTime   time.Time
	LastBorrowedUserID int
}

func NewBooksStore(db *gorm.DB) BooksStore {
	return &books{db}
}

var ErrBookNotExist = errors.New("书籍不存在")
var ErrBookHasBeenBorrowed = errors.New("书籍已借出")
var ErrBookUpdate = errors.New("书籍更新失败")

type CreateBookOptions struct {
	Name   string
	Author string
	Logo   string
	ISBN   int
}

func (db *books) Create(ctx context.Context, opts CreateBookOptions) error {
	b := &Book{
		Name:               opts.Name,
		Author:             opts.Author,
		ISBN:               opts.ISBN,
		Logo:               opts.Logo,
		IsBorrowed:         false,
		LastBorrowedTime:   time.Time{},
		LastBorrowedUserID: -1,
	}
	if err := db.WithContext(ctx).Create(b).Error; err != nil {
		log.Warn("Create book fail: %+v.%s", opts, err.Error())
		return err
	}
	log.Info("Create book: %s", opts.Name)
	return nil
}

type GetBookOptions struct {
	ID     uint
	ISBN   int
	Name   string
	Author string
}

func (db *books) Get(ctx context.Context, opts GetBookOptions) (*Book, error) {
	q := db.WithContext(ctx).Model(&Book{})

	if opts.ID != 0 {
		q = q.Where("id = ?", opts.ID)
	}
	if opts.ISBN != 0 {
		q = q.Where("isbn = ?", opts.ID)
	}
	if opts.Name != "" {
		q = q.Where("name = ?", opts.Name)
	}
	if opts.Author != "" {
		q = q.Where("author = ?", opts.Author)
	}

	var book Book
	if err := q.First(&book).Error; err != nil {
		return nil, errors.Wrap(err, "get book")
	}
	return &book, nil
}

type ListBookOption struct {
	ID  uint
	Num int
}

func (db *books) List(ctx context.Context, option ListBookOption) ([]Book, int, error) {
	b := make([]Book, 10)
	if err := db.WithContext(ctx).Where("id >= ?", option.ID).Limit(option.Num).Order("id ASC").Find(&b).Error; err != nil {
		log.Error("List book error: %v", err)
		return nil, -1, err
	}
	count := new(int64)
	if err := db.Model(Book{}).Count(count).Error; err != nil {
		return nil, -1, err
	}

	return b, int(*count), nil
}

func (db *books) Delete(ctx context.Context, id int) error {
	if err := db.WithContext(ctx).Model(Book{}).Delete("id = ?", id).Error; err != nil {
		return ErrBookNotExist
	}
	return nil
}

func (db *books) GetByID(ctx context.Context, id int) (*Book, error) {
	b := &Book{}
	if err := db.WithContext(ctx).Model(Book{}).Where("id = ?", id).First(b).Error; err != nil {
		log.Error("GetByID error: %v", err)
		return nil, err
	}
	return b, nil
}

func (db *books) GetByISBN(ctx context.Context, isbn int) ([]Book, error) {
	b := make([]Book, 10)
	if err := db.WithContext(ctx).Model(Book{}).Where("isbn = ?", isbn).First(&b).Error; err != nil {
		log.Error("GetByISBN error: %v", err)
		return nil, err
	}
	return b, nil
}

//func (db *books) UpdateByID(ctx context.Context, book Book) error {
//	return db.WithContext(ctx).Model(Book{}).Updates(book).Error
//}

func (db *books) Borrow(ctx context.Context, bid int, uid int) error {
	book, err := db.GetByID(ctx, bid)
	if err != nil {
		return ErrBookNotExist
	}

	if !book.IsBorrowed {
		if err := db.Model(&Book{}).Where("id = ?", book.ID).Updates(&Book{
			IsBorrowed:         true,
			LastBorrowedTime:   time.Now(),
			LastBorrowedUserID: uid,
		}).Error; err != nil {
			return ErrBookUpdate
		} else {
			return nil
		}

	} else {
		return ErrBookHasBeenBorrowed
	}
}
