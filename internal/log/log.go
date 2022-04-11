package log

import (
	"time"

	log "unknwon.dev/clog/v2"
)

func Init() {
	err := log.NewConsole(100,
		log.ConsoleConfig{
			Level: log.LevelInfo,
		})
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
	err = log.NewFile(100,
		log.FileConfig{
			Level:    log.LevelTrace,
			Filename: "./logs/mclib-" + time.Now().Format("0102T15_04_05") + ".log",
		})
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func Stop() {
	log.Stop()
}
