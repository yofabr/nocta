package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime),
	}
}

func (l *DefaultLogger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *DefaultLogger) Error(msg string) {
	l.errorLogger.Println(msg)
}

func (l *DefaultLogger) Debug(msg string) {
	l.debugLogger.Println(msg)
}
