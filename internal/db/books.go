package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	log "unknwon.dev/clog/v2"
)

var _ BooksStore = (*books)(nil)

var Books BooksStore

type BooksStore interface {
}

type books struct {
	*gorm.DB
}

type Book struct {
	gorm.Model

	Name             string
	Author           string
	ISBN             int
	Category         int
	IsBorrowed       bool
	LastBorrowedTime time.Time
}

func NewBooksStore(db *gorm.DB) BooksStore {
	return &books{db}
}

type CreateBookOptions struct {
	Name     string
	Author   string
	ISBN     int
	Category int
}

func (db *books) Create(ctx context.Context, opts CreateBookOptions) error {
	b := &Book{
		Name:             opts.Name,
		Author:           opts.Author,
		ISBN:             opts.ISBN,
		Category:         opts.Category,
		IsBorrowed:       false,
		LastBorrowedTime: nil,
	}
	if err := db.WithContext(ctx).Create(b).Error; err != nil {
		log.Warn("Create book fail: %+v.%s", opts, err.Error())
		return err
	}
	return nil
}
