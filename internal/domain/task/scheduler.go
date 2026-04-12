package task

import (
    "context"
    "time"
)

type Scheduler interface {
    Start(ctx context.Context) error
    Stop()
    ScheduleTask(ctx context.Context, task Task) (*Task, error)
    UnscheduleTask(ctx context.Context, taskID int64) error
}

type NextRunCalculator interface {
    Calculate(task Task, now time.Time) (time.Time, error)
    CalculateDaily(freq DailyFrequency, now time.Time) (time.Time, error)
    CalculateMonthly(freq MonthlyFrequency, now time.Time) (time.Time, error)
    CalculateSpecificDates(freq SpecificDatesFrequency, now time.Time) (time.Time, error)
    CalculateParity(freq ParityFrequency, now time.Time) (time.Time, error)
}