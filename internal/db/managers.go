package db

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"
)

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

func (u *Manager) EncodePassword() error {
	password, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		return errors.New("encode password")
	}
	u.Password = password
	return nil
}

func (u *Manager) ValidatePassword(password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		log.Error("validate password")
	}
	return match
}

type CreateManagerOption struct {
	Name     string `binding:"required"`
	Password string `binding:"required"`
}

func (db *managers) Create(ctx context.Context, opts CreateManagerOption) (*Manager, error) {
	m := &Manager{
		Name:     opts.Name,
		Password: opts.Password,
	}
	err := m.EncodePassword()
	if err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Create(m).Error; err != nil {
		log.Error("Create manager:%v", err)
		return nil, err
	}
	return m, nil
}

func (db *managers) GetByName(ctx context.Context, name string) (*Manager, error) {
	m := &Manager{}
	if err := db.WithContext(ctx).Model(&Manager{}).Where("name = ?", name).First(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (db *managers) GetByID(ctx context.Context, id uint) (*Manager, error) {
	m := &Manager{}
	if err := db.WithContext(ctx).Model(&Manager{}).Where("id = ?", id).First(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}
