package db

import "gorm.io/gorm"

type Manager struct {
	gorm.Model

	Name     string
	Password string
	Salt     string
}