package db

import "gorm.io/gorm"
import "github.com/lib/pq"

type User struct {
	gorm.Model

	Name       string
	Password   string
	Salt       string
	BooksLimit int
	BooksID    pq.Int32Array `gorm:"type:integer[]"`
}
