package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	NONE
)

var (
	logger   *log.Logger = log.New(os.Stdout, "", log.LstdFlags)
	logLevel Level       = INFO
	mu       sync.RWMutex
)

func Initialize(level string) error {
	mu.Lock()
	defer mu.Unlock()

	switch strings.ToLower(level) {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	case "none":
		logLevel = NONE
	default:
		return fmt.Errorf("unknown log level: %s", level)
	}
	return nil
}

func Debug(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= DEBUG {
		logger.Printf("[DEBUG] %s", msg)
	}
}

func Info(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= INFO {
		logger.Printf("[INFO] %s", msg)
	}
}

func Warn(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= WARN {
		logger.Printf("[WARN] %s", msg)
	}
}

func Error(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= ERROR {
		logger.Printf("[ERROR] %s", msg)
	}
}

func RequestLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		Debug(fmt.Sprintf("got incoming HTTP request method=%s path=%s", r.Method, r.URL.Path))
		h(w, r)
		duration := time.Since(start)
		Info(fmt.Sprintf("handled request method=%s path=%s duration=%v", r.Method, r.URL.Path, duration))
	})
}
