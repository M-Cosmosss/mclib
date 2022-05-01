package db

import (
	"context"
	"github.com/M-Cosmosss/mclib/internal/dbutil"
	"github.com/pkg/errors"
	"time"
	log "unknwon.dev/clog/v2"

	"gorm.io/gorm"
)

var _ BooksStore = (*books)(nil)

var Books BooksStore

type BooksStore interface {
	Create(ctx context.Context, opts CreateBookOptions) error
	GetByID(ctx context.Context, id int) (*Book, error)
	GetByISBN(ctx context.Context, isbn int) ([]Book, error)
	List(ctx context.Context, option ListBookOption) ([]Book, error)
	Borrow(ctx context.Context, id uint) error
}

type books struct {
	*gorm.DB
}

type Book struct {
	gorm.Model

	Name             string
	Author           string
	Logo             string
	ISBN             int
	IsBorrowed       bool
	LastBorrowedTime dbutil.NullTime
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
		Name:             opts.Name,
		Author:           opts.Author,
		ISBN:             opts.ISBN,
		Logo:             opts.Logo,
		IsBorrowed:       false,
		LastBorrowedTime: dbutil.NullTime{},
	}
	if err := db.WithContext(ctx).Create(b).Error; err != nil {
		log.Warn("Create book fail: %+v.%s", opts, err.Error())
		return err
	}
	log.Info("Create book: %s", opts.Name)
	return nil
}

type ListBookOption struct {
	ID  uint
	Num int
}

func (db *books) List(ctx context.Context, option ListBookOption) ([]Book, error) {
	b := make([]Book, 10)
	if err := db.WithContext(ctx).Where("id >= ?", option.ID).Limit(option.Num).Order("id ASC").Find(&b).Error; err != nil {
		log.Error("List book error: %v", err)
		return nil, err
	}
	return b, nil
}

func (db *books) GetByID(ctx context.Context, id int) (*Book, error) {
	b := &Book{}
	if err := db.WithContext(ctx).Model(Book{}).Where("id = ?", id).Find(b).Error; err != nil {
		log.Error("GetByID error: %v", err)
		return nil, err
	}
	return b, nil
}

func (db *books) GetByISBN(ctx context.Context, isbn int) ([]Book, error) {
	b := make([]Book, 10)
	if err := db.WithContext(ctx).Model(Book{}).Where("isbn = ?", isbn).Find(&b).Error; err != nil {
		log.Error("GetByISBN error: %v", err)
		return nil, err
	}
	return b, nil
}

func (db *books) Update(ctx context.Context, book Book) error {
	return db.WithContext(ctx).Save(&book).Error
}

func (db *books) Borrow(ctx context.Context, id uint) error {
	//b := &Book{}
	if book, err := db.GetByID(ctx, id); err != nil {
		return ErrBookNotExist
	} else {
		if !book.IsBorrowed {
			book.IsBorrowed = true
			book.LastBorrowedTime = dbutil.NullTime(time.Now())
			if err := db.Save(&book).Error; err != nil {
				return ErrBookUpdate
			} else {
				return nil
			}
		} else {
			return ErrBookHasBeenBorrowed
		}
	}
}
