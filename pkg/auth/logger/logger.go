package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Logger handles async logging to both console and file
type Logger struct {
	file      *os.File
	logChan   chan string
	wg        sync.WaitGroup
	mu        sync.Mutex
	isRunning bool
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			logChan:   make(chan string, 1000), // Buffer for 1000 log messages
			isRunning: false,
		}
	})
	return instance
}

// Init initializes the logger with a log file
func (l *Logger) Init(logFilePath string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isRunning {
		return nil // Already initialized
	}

	// Create or open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.file = file
	l.isRunning = true

	// Start async log writer
	l.wg.Add(1)
	go l.logWriter()

	return nil
}

// logWriter processes log messages asynchronously
func (l *Logger) logWriter() {
	defer l.wg.Done()

	for msg := range l.logChan {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMsg := fmt.Sprintf("[%s] %s\n", timestamp, msg)

		// Write to console (stdout)
		fmt.Print(logMsg)

		// Write to file
		if l.file != nil {
			l.file.WriteString(logMsg)
		}
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf("[INFO] "+format, args...)
	l.log(msg)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf("[ERROR] "+format, args...)
	l.log(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf("[WARN] "+format, args...)
	l.log(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf("[DEBUG] "+format, args...)
	l.log(msg)
}

// log sends a message to the log channel
func (l *Logger) log(msg string) {
	if !l.isRunning {
		// Fallback to standard log if logger not initialized
		log.Println(msg)
		return
	}

	select {
	case l.logChan <- msg:
		// Message sent successfully
	default:
		// Channel full, log to stderr
		fmt.Fprintf(os.Stderr, "Log channel full, dropping message: %s\n", msg)
	}
}

// Close closes the logger and waits for all logs to be written
func (l *Logger) Close() {
	l.mu.Lock()
	if !l.isRunning {
		l.mu.Unlock()
		return
	}
	l.isRunning = false
	l.mu.Unlock()

	// Close the channel and wait for all logs to be written
	close(l.logChan)
	l.wg.Wait()

	// Close the file
	if l.file != nil {
		l.file.Close()
	}
}

// Helper functions for global logger
func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

