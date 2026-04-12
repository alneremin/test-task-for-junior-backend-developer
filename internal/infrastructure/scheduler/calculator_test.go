package scheduler

import (
    "testing"
    "time"
    
    taskdomain "example.com/taskservice/internal/domain/task"
)

func TestNextRunCalculatorImpl_Calculate(t *testing.T) {
    calculator := NewNextRunCalculator()
    now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
    
    tests := []struct {
        name      string
        task      taskdomain.Task
        now       time.Time
        want      time.Time
        wantError error
    }{
        {
            name: "daily frequency",
            task: taskdomain.Task{
                Frequency: taskdomain.FrequencyRule{
                    Type:  taskdomain.Daily,
                    Daily: &taskdomain.DailyFrequency{Hour: 14, Minute: 30, Interval: 1},
                },
            },
            now:  now,
            want: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "monthly frequency",
            task: taskdomain.Task{
                Frequency: taskdomain.FrequencyRule{
                    Type:    taskdomain.Monthly,
                    Monthly: &taskdomain.MonthlyFrequency{DayOfMonth: 20, Hour: 9, Minute: 0},
                },
            },
            now:  now,
            want: time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
        },
        {
            name: "unknown frequency type",
            task: taskdomain.Task{
                Frequency: taskdomain.FrequencyRule{
                    Type: "unknown",
                },
            },
            now:       now,
            wantError: taskdomain.ErrUnknownFrequencyType,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.Calculate(tt.task, tt.now)
            if tt.wantError != nil {
                if err != tt.wantError {
                    t.Errorf("Calculate() error = %v, want %v", err, tt.wantError)
                }
                return
            }
            if err != nil {
                t.Errorf("Calculate() unexpected error: %v", err)
            }
            if !got.Equal(tt.want) {
                t.Errorf("Calculate() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestNextRunCalculatorImpl_CalculateDaily(t *testing.T) {
    calculator := &NextRunCalculatorImpl{}
    
    tests := []struct {
        name     string
        freq     taskdomain.DailyFrequency
        fromTime time.Time
        want     time.Time
    }{
        {
            name:     "same day - future time",
            freq:     taskdomain.DailyFrequency{Hour: 14, Minute: 30, Interval: 1},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
        },
        {
            name:     "same day - past time, next day",
            freq:     taskdomain.DailyFrequency{Hour: 8, Minute: 0, Interval: 1},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 16, 8, 0, 0, 0, time.UTC),
        },
        {
            name:     "interval 3 days",
            freq:     taskdomain.DailyFrequency{Hour: 14, Minute: 30, Interval: 3},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 17, 14, 30, 0, 0, time.UTC),
        },
        {
            name:     "interval 1 day - exact time match",
            freq:     taskdomain.DailyFrequency{Hour: 10, Minute: 0, Interval: 1},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
        },
        {
            name:     "end of month transition",
            freq:     taskdomain.DailyFrequency{Hour: 23, Minute: 59, Interval: 1},
            fromTime: time.Date(2024, 1, 31, 22, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 31, 23, 59, 0, 0, time.UTC),
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.CalculateDaily(tt.freq, tt.fromTime)
            if err != nil {
                t.Errorf("CalculateDaily() unexpected error: %v", err)
            }
            if !got.Equal(tt.want) {
                t.Errorf("CalculateDaily() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestNextRunCalculatorImpl_CalculateMonthly(t *testing.T) {
    calculator := &NextRunCalculatorImpl{}
    
    tests := []struct {
        name     string
        freq     taskdomain.MonthlyFrequency
        fromTime time.Time
        want     time.Time
    }{
        {
            name:     "same month - future date",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 20, Hour: 14, Minute: 30},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 20, 14, 30, 0, 0, time.UTC),
        },
        {
            name:     "same month - past date, next month",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 10, Hour: 14, Minute: 30},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 2, 10, 14, 30, 0, 0, time.UTC),
        },
        {
            name:     "february - day 30 adjusted to 28",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 30, Hour: 12, Minute: 0},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 30, 12, 0, 0, 0, time.UTC),
        },
        {
            name:     "february in leap year - day 30 adjusted to 29",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 30, Hour: 12, Minute: 0},
            fromTime: time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
        },
        {
            name:     "december to january transition",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 5, Hour: 9, Minute: 0},
            fromTime: time.Date(2024, 12, 10, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2025, 1, 5, 9, 0, 0, 0, time.UTC),
        },
        {
            name:     "exact time match - same day and time",
            freq:     taskdomain.MonthlyFrequency{DayOfMonth: 15, Hour: 10, Minute: 0},
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.CalculateMonthly(tt.freq, tt.fromTime)
            if err != nil {
                t.Errorf("CalculateMonthly() unexpected error: %v", err)
            }
            if !got.Equal(tt.want) {
                t.Errorf("CalculateMonthly() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestNextRunCalculatorImpl_CalculateSpecificDates(t *testing.T) {
    calculator := &NextRunCalculatorImpl{}
    
    tests := []struct {
        name      string
        freq      taskdomain.SpecificDatesFrequency
        fromTime  time.Time
        want      time.Time
        wantError bool
    }{
        {
            name: "single future date",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024-01-20"},
                Hour:  14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 20, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "multiple dates - pick earliest future",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024-01-25", "2024-01-20", "2024-01-30"},
                Hour:  9, Minute: 0,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
        },
        {
            name: "today's date with future time",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024-01-15"},
                Hour:  14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            want:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "today's date with past time - include anyway",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024-01-15"},
                Hour:  8, Minute: 0,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantError: true,
        },
        {
            name: "all dates in past",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024-01-10", "2024-01-12"},
                Hour:  14, Minute: 30,
            },
            fromTime:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantError: true,
        },
        {
            name: "invalid date format",
            freq: taskdomain.SpecificDatesFrequency{
                Dates: []string{"2024/01/20"},
                Hour:  14, Minute: 30,
            },
            fromTime:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
            wantError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.CalculateSpecificDates(tt.freq, tt.fromTime)
            if tt.wantError {
                if err == nil {
                    t.Errorf("CalculateSpecificDates() expected error but got none")
                }
                return
            }
            if err != nil {
                t.Errorf("CalculateSpecificDates() unexpected error: %v", err)
            }
            if !got.Equal(tt.want) {
                t.Errorf("CalculateSpecificDates() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestNextRunCalculatorImpl_CalculateParity(t *testing.T) {
    calculator := &NextRunCalculatorImpl{}
    
    tests := []struct {
        name     string
        freq     taskdomain.ParityFrequency
        fromTime time.Time
        want     time.Time
    }{
        {
            name: "even day - same day",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Even,
                Hour: 14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), // 16 - even
            want:     time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "even day - need next day",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Even,
                Hour: 14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), // 15 - odd, next even is 16
            want:     time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "odd day - same day",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Odd,
                Hour: 14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), // 15 - odd
            want:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "odd day - need next day",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Odd,
                Hour: 14, Minute: 30,
            },
            fromTime: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), // 16 - even, next odd is 17
            want:     time.Date(2024, 1, 17, 14, 30, 0, 0, time.UTC),
        },
        {
            name: "even day - time passed today",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Even,
                Hour: 8, Minute: 0,
            },
            fromTime: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), // time passed, but day is even
            want:     time.Date(2024, 1, 18, 8, 0, 0, 0, time.UTC),  // next even day
        },
        {
            name: "month transition",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Even,
                Hour: 12, Minute: 0,
            },
            fromTime: time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC), // 31 - odd, next even is Feb 2
            want:     time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC),
        },
        {
            name: "year transition",
            freq: taskdomain.ParityFrequency{
                Type: taskdomain.Odd,
                Hour: 12, Minute: 0,
            },
            fromTime: time.Date(2024, 12, 30, 10, 0, 0, 0, time.UTC), // 30 - even, next odd is Dec 31
            want:     time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC),
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.CalculateParity(tt.freq, tt.fromTime)
            if err != nil {
                t.Errorf("CalculateParity() unexpected error: %v", err)
            }
            if !got.Equal(tt.want) {
                t.Errorf("CalculateParity() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Бенчмарк тесты
func BenchmarkCalculateDaily(b *testing.B) {
    calculator := &NextRunCalculatorImpl{}
    freq := taskdomain.DailyFrequency{Hour: 14, Minute: 30, Interval: 1}
    fromTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
    
    for i := 0; i < b.N; i++ {
        _, _ = calculator.CalculateDaily(freq, fromTime)
    }
}

func BenchmarkCalculateMonthly(b *testing.B) {
    calculator := &NextRunCalculatorImpl{}
    freq := taskdomain.MonthlyFrequency{DayOfMonth: 20, Hour: 14, Minute: 30}
    fromTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
    
    for i := 0; i < b.N; i++ {
        _, _ = calculator.CalculateMonthly(freq, fromTime)
    }
}