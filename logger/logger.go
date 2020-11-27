package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	levelDebug = iota
	levelError
	levelFatal
)

const (
	printDebug = "[__Debug__]"
	printError = "[__Error__]"
	printFatal = "[__Fatal__]"
)

type Logger struct {
	logLevel int
	baseLogger *log.Logger
	logFile *os.File
}

func New(logLevel int) (*Logger, error){
	now := time.Now()
	fileName := fmt.Sprintf("log/%v_%v/%v_%v:%v", now.Year(), now.Month(), now.Day(), now.Hour(), now.Second())
	logFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	baseLogger := log.New(logFile, "", log.LUTC)

	return &Logger{
		logLevel: logLevel,
		baseLogger: baseLogger,
		logFile: logFile,
	}, nil
}

func (logger *Logger) Print(level int, printLevel string, format string, v ...interface{}) {
	if level < logger.logLevel {
		return
	}
	format = printLevel + format + "\n"
	logger.baseLogger.Printf(format, v...)

	if level == levelFatal {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, v ...interface{}) {
	logger.Print(levelDebug, printDebug, format, v...)
}

func (logger *Logger) Error(format string, v ...interface{}) {
	logger.Print(levelError, printError, format, v...)
}

func (logger *Logger) Fatal(format string, v ...interface{}) {
	logger.Print(levelFatal, printFatal, format, v...)
}

var gLoger *Logger

func init() {
	gLoger, _ = New(levelDebug)
}

func Debug (format string, v ...interface{}) {
	gLoger.Debug(format, v...)
}

func Error (format string, v ...interface{}) {
	gLoger.Error(format, v...)
}

func Fatal (format string, v ...interface{}) {
	gLoger.Fatal(format, v...)
}