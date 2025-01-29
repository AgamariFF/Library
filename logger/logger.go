package logger

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	logFile  *os.File
)

func InitLog() error {
	f, err := os.OpenFile("/app/logger/info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	InfoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func CloseLog() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
