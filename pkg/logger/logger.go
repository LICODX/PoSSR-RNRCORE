package logger

import (
	"log"
	"os"
)

// Logger provides structured logging
type Logger struct {
	level Level
}

// Level represents log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var globalLogger = &Logger{level: LevelInfo}

// SetLevel sets global log level
func SetLevel(level Level) {
	globalLogger.level = level
}

// Info logs info message
func Info(format string, args ...interface{}) {
	if globalLogger.level <= LevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

// Debug logs debug message
func Debug(format string, args ...interface{}) {
	if globalLogger.level <= LevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Warn logs warning
func Warn(format string, args ...interface{}) {
	if globalLogger.level <= LevelWarn {
		log.Printf("[WARN] "+format, args...)
	}
}

// Error logs error
func Error(format string, args ...interface{}) {
	if globalLogger.level <= LevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

// Fatal logs fatal error and exits
func Fatal(format string, args ...interface{}) {
	log.Printf("[FATAL] "+format, args...)
	os.Exit(1)
}

// Setup initializes logger with file output
func Setup(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}
