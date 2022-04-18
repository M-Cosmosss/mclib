package db

import "gorm.io/gorm"

var _ ManagersStore = (*managers)(nil)

var Managers ManagersStore

type ManagersStore interface {
}

type managers struct {
	*gorm.DB
}

type Manager struct {
	gorm.Model

	Name     string
	Password string
	Salt     string
}

func NewManagersStore(db *gorm.DB) ManagersStore {
	return &managers{db}
}
