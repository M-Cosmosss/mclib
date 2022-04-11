package db

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	gorm.Model

	Name             string
	Author           string
	Pic              string
	ISBN             int
	Category         int
	IsBorrowed       bool
	LastBorrowedTime time.Time
}
