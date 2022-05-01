package db

import (
	"context"
	"github.com/M-Cosmosss/mclib/internal/config"
	"github.com/alexedwards/argon2id"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	log "unknwon.dev/clog/v2"
)

var _ UsersStore = (*users)(nil)

var Users UsersStore

type UsersStore interface {
	Create(ctx context.Context, opts CreateUserOption) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Borrow(ctx context.Context, user User, id int) error
	Return(ctx context.Context, user User, nums []int32) error
	DeleteByID(ctx context.Context, id uint) error
}

type users struct {
	*gorm.DB
}

type User struct {
	gorm.Model

	Name       string `gorm:"unique"`
	Password   string
	BooksLimit int
	IsManager  bool
	BooksID    pq.Int32Array `gorm:"type:integer[]"`
}

func (u *User) EncodePassword() error {
	password, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		return errors.New("encode password")
	}
	u.Password = password
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		log.Error("validate password")
	}
	return match
}

func NewUsersStore(db *gorm.DB) UsersStore {
	return &users{db}
}

var ErrUserNotExists = errors.New("用户不存在")

type CreateUserOption struct {
	Name     string `binding:"required"`
	Password string `binding:"required"`
	IsAdmin  bool
}

func (db *users) Create(ctx context.Context, opts CreateUserOption) (*User, error) {
	u := &User{
		Name:       opts.Name,
		Password:   opts.Password,
		IsManager:  opts.IsAdmin,
		BooksLimit: config.DefaultUserBooksLimit,
	}
	err := u.EncodePassword()
	if err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		log.Error("Create user:%v", err)
		return nil, err
	}
	return u, nil
}

func (db *users) GetByName(ctx context.Context, name string) (*User, error) {
	u := &User{}
	if err := db.WithContext(ctx).Model(&User{}).Where("name = ?", name).First(u).Error; err != nil {
		return nil, ErrUserNotExists
	}
	return u, nil
}

func (db *users) GetByID(ctx context.Context, id uint) (*User, error) {
	u := &User{}
	if err := db.WithContext(ctx).Model(&User{}).Where("id = ?", id).First(u).Error; err != nil {
		return nil, ErrUserNotExists
	}
	if u.BooksID == nil {
		u.BooksID = []int32{}
	}
	return u, nil
}

func (db *users) Borrow(ctx context.Context, user User, id int) error {
	user.BooksID = append(user.BooksID, int32(id))
	if err := db.WithContext(ctx).Model(&User{}).
		Where("id = ?", user.ID).
		Update("books_id", user.BooksID).Error; err != nil {
		return err
	}
	return nil
}

func (db *users) Return(ctx context.Context, user User, nums []int32) error {
	user.BooksID = nums
	if err := db.WithContext(ctx).Model(User{}).Where("id = ?", user.ID).Update("books_id", user.BooksID).Error; err != nil {
		return err
	}
	return nil
}

func (db *users) DeleteByID(ctx context.Context, id uint) error {
	_, err := db.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get by ID")
	}
	return db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}
