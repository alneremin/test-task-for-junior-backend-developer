package task

import "errors"

var ErrTitleIsEmpty = errors.New("title is required")
var ErrInvalidStatus = errors.New("invalid status")
var ErrIdMustBetPositive = errors.New("id must be positive")