package db

import (
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

func Init() {
	conf := readConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.DBName)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info("PostgreSQL connected.")
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
