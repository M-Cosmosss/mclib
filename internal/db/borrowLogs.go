package db

import "gorm.io/gorm"

const (
	BORROW = iota
	RETURN
)

type BorrowLog struct {
	gorm.Model

	UserID    int
	BookID    int
	Operation int
}
