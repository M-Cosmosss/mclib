package db

import "gorm.io/gorm"

const (
	BORROW = iota
	RETURN
)

var _ BorrowLogsStore = (*borrowLogs)(nil)

var BorrowLogs BorrowLogsStore

type BorrowLogsStore interface {
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
