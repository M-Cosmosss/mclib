package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	log "unknwon.dev/clog/v2"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pgConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	Host     string `json:"host"`
}

var AllTables = []interface{}{
	&Book{},
	&User{},
	&Manager{},
	&BorrowLog{},
}

type Transactor interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
}

func Init() {
	conf := readConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("PostgreSQL connect error: %v", err)
	}
	log.Info("PostgreSQL connected.")
	if err = db.AutoMigrate(AllTables...); err != nil {
		log.Fatal("AutoMigrate error: %v", err)
	}

	SetDatabaseStore(db)

	Users.Create(context.Background(), CreateUserOption{
		Name:     "admin",
		Password: "admin",
		IsAdmin:  true,
	})
}

func readConfig() *pgConfig {
	conf := &pgConfig{}
	b, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func SetDatabaseStore(db *gorm.DB) {
	Books = NewBooksStore(db)
	BorrowLogs = NewBorrowLogsStore(db)
	Managers = NewManagersStore(db)
	Users = NewUsersStore(db)
}
