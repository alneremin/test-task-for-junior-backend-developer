package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            			`json:"title"`
	Description string           		 	`json:"description"`
	Status      taskdomain.Status 			`json:"status"`
	Frequency   taskdomain.FrequencyRule 	`json:"frequency"`
}

type taskDTO struct {
	ID          int64           		 	`json:"id"`
	Title       string           		 	`json:"title"`
	Description string          		 	`json:"description"`
	Status      taskdomain.Status 			`json:"status"`
	Frequency   taskdomain.FrequencyRule 	`json:"frequency"`
	NextRunTime time.Time   				`json:"next_run_time"`
	CreatedAt   time.Time         			`json:"created_at"`
	UpdatedAt   time.Time        		 	`json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Frequency:   task.Frequency,
		NextRunTime: task.NextRunTime,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
