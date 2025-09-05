package config

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr   string
	FlagLogLevel  string
	FlagWorkers   int
	FlagQueueSize int
)

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&FlagWorkers, "workers", 4, "number of worker goroutines")
	flag.IntVar(&FlagQueueSize, "queue-size", 64, "buffered queue size")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
	if v := os.Getenv("WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			FlagWorkers = n
		}
	}
	if v := os.Getenv("QUEUE_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			FlagQueueSize = n
		}
	}
}
