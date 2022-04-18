package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var _ UsersStore = (*users)(nil)

var Users UsersStore

type UsersStore interface {
}

type users struct {
	*gorm.DB
}

type User struct {
	gorm.Model

	Name       string
	Password   string
	Salt       string
	BooksLimit int
	BooksID    pq.Int32Array `gorm:"type:integer[]"`
}

func NewUsersStore(db *gorm.DB) UsersStore {
	return &users{db}
}

func (db *users) New() {

}
