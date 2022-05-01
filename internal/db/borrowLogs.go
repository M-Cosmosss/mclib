package db

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

const (
	BORROW = iota
	RETURN
)

var _ BorrowLogsStore = (*borrowLogs)(nil)

var BorrowLogs BorrowLogsStore

type BorrowLogsStore interface {
	List(ctx context.Context, opts ListBorrowLogsOptions) ([]*BorrowLog, int64, error)
	Create(ctx context.Context, opt CreateLogOption) error
	GetByBookID(ctx context.Context, option GetLogByBookIDOption) ([]BorrowLog, error)
}

type borrowLogs struct {
	*gorm.DB
}

type BorrowLog struct {
	gorm.Model

	UserID    int
	BookID    int
	Operation int
}

func NewBorrowLogsStore(db *gorm.DB) BorrowLogsStore {
	return &borrowLogs{db}
}

type ListBorrowLogsOptions struct {
	Limit  int
	Offset int
}

func (db *borrowLogs) List(ctx context.Context, opts ListBorrowLogsOptions) ([]*BorrowLog, int64, error) {
	q := db.WithContext(ctx).Model(&BorrowLog{})

	// 查询总数
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return nil, count, nil
	}

	var borrowLogs []*BorrowLog
	if err := q.Limit(opts.Limit).Offset(opts.Offset).Find(&borrowLogs).Error; err != nil {
		return nil, 0, errors.Wrap(err, "find")
	}
	return borrowLogs, count, nil
}

type CreateLogOption struct {
	User      int
	Book      int
	Operation int
}

func (db *borrowLogs) Create(ctx context.Context, opt CreateLogOption) error {
	l := &BorrowLog{
		UserID:    opt.User,
		BookID:    opt.Book,
		Operation: opt.Operation,
	}
	if err := db.WithContext(ctx).Create(l).Error; err != nil {
		log.Error("Create log error: %v")
		return err
	}
	return nil
}

type GetLogByBookIDOption struct {
	BookID int
	Offset int
	Limit  int
}

func (db *borrowLogs) GetByBookID(ctx context.Context, option GetLogByBookIDOption) ([]BorrowLog, error) {
	var re []BorrowLog
	if err := db.WithContext(ctx).Where("book_id = ?", option.BookID).Order("id DESC").Limit(option.Limit).Offset(option.Offset).Find(&re).Error; err != nil {
		return nil, err
	}
	return re, nil
}
