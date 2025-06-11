package main

import (
	"fmt"
	"log"
)

// Logger provides structured logging for the OGP generator.
type Logger struct{}

// NewLogger creates a new Logger instance.
func NewLogger() *Logger {
	return &Logger{}
}

// Warning logs a warning message.
func (l *Logger) Warning(format string, args ...interface{}) {
	fmt.Printf("Warning: "+format+"\n", args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
}

// Info logs an informational message.
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf("Info: "+format+"\n", args...)
}

// Fatal logs a fatal error and exits the program.
func (l *Logger) Fatal(format string, args ...interface{}) {
	log.Fatalf("Fatal: "+format, args...)
}

// DefaultLogger is the global logger instance used throughout the application.
var DefaultLogger = NewLogger()
