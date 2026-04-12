package task

import (
	"context"
	"fmt"
	"time"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo 		taskdomain.Repository
	scheduler  	taskdomain.Scheduler
	calculator  taskdomain.NextRunCalculator
	now  func() time.Time
}

func NewService(repo taskdomain.Repository, scheduler taskdomain.Scheduler, calculator taskdomain.NextRunCalculator) *Service {
	return &Service{
		repo: repo,
		scheduler: scheduler,
		calculator: calculator,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input taskdomain.CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Frequency:   normalized.Frequency,
	}
	now := s.now()
	if !normalized.Frequency.IsEmpty() {
		next_time, err := s.calculator.Calculate(*model, now)
		if err != nil {
			return nil, err
		}
		model.NextRunTime = next_time
	}

	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrIdMustBetPositive)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input taskdomain.UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrIdMustBetPositive)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Frequency:   normalized.Frequency,
	}
	now := s.now()
	if !normalized.Frequency.IsEmpty() {
		next_time, err := s.calculator.Calculate(*model, now)
		if err != nil {
			return nil, err
		}
		model.NextRunTime = next_time
	}
	model.UpdatedAt = now

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrIdMustBetPositive)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}
