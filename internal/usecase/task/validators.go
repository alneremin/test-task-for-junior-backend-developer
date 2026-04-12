package task

import (
	"fmt"
	"strings"
	taskdomain "example.com/taskservice/internal/domain/task"
)

func validateCreateInput(input taskdomain.CreateInput) (taskdomain.CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return taskdomain.CreateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrTitleIsEmpty)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return taskdomain.CreateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrInvalidStatus)
	}

	if isValid, errMsg := input.Frequency.Valid(); !isValid {
        return taskdomain.CreateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, errMsg)
    }

	return input, nil
}

func validateUpdateInput(input taskdomain.UpdateInput) (taskdomain.UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return taskdomain.UpdateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrTitleIsEmpty)
	}

	if !input.Status.Valid() {
		return taskdomain.UpdateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, ErrInvalidStatus)
	}

	if isValid, errMsg := input.Frequency.Valid(); !isValid {
        return taskdomain.UpdateInput{}, fmt.Errorf("%w: %w", taskdomain.ErrInvalidInput, errMsg)
    }

	return input, nil
}
