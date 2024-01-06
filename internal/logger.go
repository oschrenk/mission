package internal

import (
	"log"
	"os"
)

type Log struct {
	Enabled bool
}

var errLog = log.New(os.Stderr, "", log.LstdFlags)

var Logger = Log{
	Enabled: false,
}

func (logger *Log) Log(s string) {
	if logger.Enabled {
		errLog.Println(s)
	}
}
