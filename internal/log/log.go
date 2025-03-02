package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	generalLogger *log.Logger
	focusLogger   *log.Logger
	once          sync.Once
)

const (
	logFlags = log.LstdFlags | log.Lshortfile
)

func init() {
	once.Do(initLoggers)
}

func initLoggers() {
	// Create or open log files
	generalFile, err := os.OpenFile("internal/log/general.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open general log file: %v", err)
	}

	focusFile, err := os.OpenFile("internal/log/focus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open focus log file: %v", err)
	}

	// Initialize loggers
	generalLogger = log.New(io.MultiWriter(os.Stdout, generalFile), "INFO: ", logFlags)
	focusLogger = log.New(io.MultiWriter(os.Stdout, focusFile), "FOCUS: ", logFlags)
}

// Info logs general information
func Info(message string) {
	generalLogger.SetPrefix("INFO: ")
	generalLogger.Output(2, message)
}

// Debug logs debug information
func Debug(message string) {
	generalLogger.SetPrefix("DEBUG: ")
	generalLogger.Output(2, message)
}

// Warn logs warning messages
func Warn(message string) {
	focusLogger.SetPrefix("WARNING: ")
	focusLogger.Output(2, message)
}

// Critical logs critical errors
func Critical(message string) {
	focusLogger.SetPrefix("CRITICAL: ")
	focusLogger.Output(2, message)
}
