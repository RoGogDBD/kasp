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

	levelStr := strings.ToLower(strings.TrimSpace(level))
	switch levelStr {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
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
		logger.Printf("\033[32m[DEBUG]\033[0m %s", msg)
	}
}

func Info(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= INFO {
		logger.Printf("\033[34m[INFO]\033[0m %s", msg)
	}
}

func Error(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	if logLevel <= ERROR {
		logger.Printf("\033[31m[ERROR]\033[0m %s", msg)
	}
}

func RequestLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		Debug(fmt.Sprintf("METHOD=%s PATH=%s", r.Method, r.URL.Path))
		h(w, r)
		duration := time.Since(start)
		Info(fmt.Sprintf("METHOD=%s PATH=%s DURATION=%v", r.Method, r.URL.Path, duration))
	})
}
