package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/RoGogDBD/kasp/internal/logger"
	"github.com/RoGogDBD/kasp/internal/models"
	"github.com/RoGogDBD/kasp/internal/repository"
)

func StartWorkers(ctx context.Context, n int, q *repository.Queue, s *repository.Storage) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, q, s)
		}(i)
	}
	go func() {
		<-ctx.Done()
		q.Close()
	}()
	wg.Wait()
}

func worker(ctx context.Context, id int, q *repository.Queue, s *repository.Storage) {
	for {
		select {
		case <-ctx.Done():
			logger.Info(fmt.Sprintf("worker %d shutting down", id))
			return
		case task, ok := <-q.Tasks():
			if !ok {
				logger.Info(fmt.Sprintf("worker %d: queue closed", id))
				return
			}
			s.SetStatus(task.ID, models.StatusRunning)
			logger.Info(fmt.Sprintf("worker %d started processing task id=%s", id, task.ID))
			processTask(task, s)
			logger.Info(fmt.Sprintf("worker %d finished task id=%s", id, task.ID))
		}
	}
}

func processTask(t *models.Task, s *repository.Storage) {
	retries := 0
	for {
		s.SetStatus(t.ID, models.StatusRunning)

		time.Sleep(time.Duration(100+rand.Intn(400)) * time.Millisecond)

		if rand.Intn(100) < 20 {
			retries++
			if retries > t.MaxRetries {
				s.SetStatus(t.ID, models.StatusFailed)
				logger.Error(fmt.Sprintf("task %s failed after %d retries", t.ID, retries-1))
				return
			}

			backoff := time.Duration(100*(1<<retries)) * time.Millisecond
			jitter := time.Duration(rand.Intn(100)) * time.Millisecond
			sleepTime := backoff + jitter
			logger.Debug(fmt.Sprintf("task %s failed, retrying in %v (attempt %d)", t.ID, sleepTime, retries))
			time.Sleep(sleepTime)
			continue
		}

		s.SetStatus(t.ID, models.StatusDone)
		return
	}
}
