package logger

import (
	"os"

	"github.com/Sirupsen/logrus"

	"constants"
)

type Logger struct {
	*logrus.Logger
}

var (
	ZaloraStatsLogger *Logger
)

func Init() {
	pwd, _ := os.Getwd()
	ZaloraStatsLogger = NewLogger(pwd + constants.LoggerFilePath)
}

func NewLogger(filePath string) *Logger {
	logger := Logger{Logger: logrus.New()}
	logger.SetFormatter(&logrus.JSONFormatter{})
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalln("cant open zalora log, ", err)
	}
	logger.Logger.Out = f
	return &logger
}

func (l *Logger) Error(bucket string, identifier string, message string, errorMessage string) {
	l.Logger.WithFields(logrus.Fields{
		"bucket":     bucket,
		"identifier": identifier,
	}).Error(message + ": " + errorMessage)
}
