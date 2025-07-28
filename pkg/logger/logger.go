package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	mu       sync.Mutex
	logLevel = INFO
	logger   = log.New(os.Stdout, "", 0)
)

// SetLevel sets the global log level (DEBUG, INFO, WARN, ERROR)
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()
	logLevel = level
}

// SetOutput allows setting a custom output (e.g., file, os.Stderr)
func SetOutput(output *os.File) {
	mu.Lock()
	defer mu.Unlock()
	logger.SetOutput(output)
}

// logMessage formats and outputs a log message with level and timestamp
func logMessage(level Level, levelStr string, format string, args ...interface{}) {
	if level < logLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	logger.Printf("[%s] [%s] %s\n", timestamp, levelStr, msg)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	logMessage(DEBUG, "DEBUG", format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	logMessage(INFO, "INFO", format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	logMessage(WARN, "WARN", format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	logMessage(ERROR, "ERROR", format, args...)
}
