package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RoGogDBD/kasp/internal/logger"
	"github.com/RoGogDBD/kasp/internal/models"
	"github.com/RoGogDBD/kasp/internal/repository"
)

func EnqueueHandler(q *repository.Queue, s *repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid payload"}`))
			return
		}

		if !q.Enqueue(&task) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"queue is full"}`))
			return
		}

		s.SetStatus(task.ID, models.StatusQueued)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"text":"Task queued"}}`))
		logger.Debug("task added to queue")
	}
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`OK`))
	}
}
