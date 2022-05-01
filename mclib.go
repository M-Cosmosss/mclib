package main

import (
	"github.com/M-Cosmosss/mclib/internal/db"
	"github.com/M-Cosmosss/mclib/internal/log"
	"github.com/M-Cosmosss/mclib/internal/route"
)

func main() {
	log.Init()
	defer log.Stop()
	db.Init()
	route.Init().Run()
}
