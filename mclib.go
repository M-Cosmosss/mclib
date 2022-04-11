package main

import (
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/M-Cosmosss/mclib/internal/log"
)

func main() {
	log.Init()
	defer log.Stop()
	db.Init()
}
