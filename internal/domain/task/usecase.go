package task

import (
	"context"
)

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*Task, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]Task, error)
}


type CreateInput struct {
	Title       string
	Description string
	Status      Status
	Frequency   FrequencyRule
}

type UpdateInput struct {
	Title       string
	Description string
	Status      Status
	Frequency   FrequencyRule
}
