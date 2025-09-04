package main

import (
	"fmt"
	"net/http"

	"github.com/RoGogDBD/kasp/internal/logger"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		logger.Error(fmt.Sprintf("server error: %v", err))
	}
}

func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/enqueue", logger.RequestLogger(webhook))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	logger.Info(fmt.Sprintf("Running server on address %s", flagRunAddr))
	return http.ListenAndServe(flagRunAddr, mux)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Debug(fmt.Sprintf("got request with bad method=%s", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`
      {
        "response": {
          "text": "Test"
        }
      }
    `))
	logger.Debug("sending HTTP 200 response")
}
