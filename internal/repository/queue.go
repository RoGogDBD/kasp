package repository

import (
	"github.com/RoGogDBD/kasp/internal/models"
)

type Queue struct {
	tasks chan *models.Task
}

func NewQueue(size int) *Queue {
	return &Queue{tasks: make(chan *models.Task, size)}
}

func (q *Queue) Enqueue(task *models.Task) bool {
	select {
	case q.tasks <- task:
		return true
	default:
		return false
	}
}

func (q *Queue) Dequeue() (*models.Task, bool) {
	select {
	case task := <-q.tasks:
		return task, true
	default:
		return nil, false
	}
}

func (q *Queue) Tasks() <-chan *models.Task {
	return q.tasks
}

func (q *Queue) Close() {
	close(q.tasks)
}
