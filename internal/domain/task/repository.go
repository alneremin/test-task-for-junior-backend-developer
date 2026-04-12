package task

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Delete(ctx context.Context, id int64) (error)
	List(ctx context.Context) ([]Task, error)
	GetTasksForExecution(ctx context.Context) ([]Task, error)
	UpdateNextRunTime(ctx context.Context, id int64, nextRun time.Time) (*Task, error)
}