package task

import (
	"errors"
	"fmt"
)

var (

	ErrInvalidInput = errors.New("invalid task input")
	ErrNotFound = errors.New("task not found")
	ErrDailyConfigIsRequired = errors.New("daily config is required")
	ErrInvalidDailyInterval = errors.New("daily interval must be between 1 and 365")
	ErrMonthlyConfigIsRequired = errors.New("monthly config is required")
	ErrInvalidMonthlyInterval = errors.New("day of month must be between 1 and 30")
	ErrSpecificDatesIsRequired = errors.New("specific dates config is required")
	ErrAtLeastOneDateIsRequired = errors.New("at least one specific date is required")
	ErrInvalidDate = errors.New("invalid date format")
	ErrInvalidHour = errors.New("invalid hour format")
	ErrInvalidMinute = errors.New("invalid minute format")
)

type InvalidDateError struct {
    Format 	string
	Date   	string
    Err    	error
}

func (e *InvalidDateError) Error() string {
    return fmt.Sprintf("%v: %s, required: %s", e.Err, e.Date, e.Format)
}

func (e *InvalidDateError) Unwrap() error {
    return e.Err
}

var (
	ErrParityConfigIsRequired = errors.New("parity config is required")
	ErrInvalidParityType = errors.New("parity type must be 'even' or 'odd'")
	ErrUnknownFrequencyType = errors.New("unknown frequency type")
)

type UnknownFrequencyTypeError struct {
    Type 	string
    Err    	error
}

func (e *UnknownFrequencyTypeError) Error() string {
    return fmt.Sprintf("%v: '%s'", e.Err, e.Type)
}

func (e *UnknownFrequencyTypeError) Unwrap() error {
    return e.Err
}
