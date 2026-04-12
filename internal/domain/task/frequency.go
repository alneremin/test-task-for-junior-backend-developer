package task

import (
    "time"
    "reflect"
)

type FrequencyType string

const (
    Daily           FrequencyType = "daily"
    Monthly         FrequencyType = "monthly"
    SpecificDates   FrequencyType = "specific_dates"
    Parity          FrequencyType = "parity"
)

type ParityType string

const (
    Even ParityType = "even"
    Odd  ParityType = "odd"
)

type DailyFrequency struct { // TODO: учесть высокосный год
    Interval int `json:"interval" validate:"min=1,max=365"` // каждый N-й день
    Hour   int `json:"hour" validate:"min=0,max=23"`
    Minute int `json:"minute" validate:"min=0,max=59"`
}

type MonthlyFrequency struct {
    DayOfMonth int `json:"day_of_month" validate:"min=1,max=31"` // число месяца
    Hour       int `json:"hour" validate:"min=0,max=23"`
    Minute     int `json:"minute" validate:"min=0,max=59"`
}

type SpecificDatesFrequency struct {
    Dates []string          `json:"dates" validate:"required,min=1"` // даты в формате YYYY-MM-DD
    Hour       int          `json:"hour" validate:"min=0,max=23"`
    Minute     int          `json:"minute" validate:"min=0,max=59"`
}

type ParityFrequency struct {
    Type       ParityType   `json:"type" validate:"oneof=even odd"`
    Hour       int          `json:"hour" validate:"min=0,max=23"`
    Minute     int          `json:"minute" validate:"min=0,max=59"`
}

// FrequencyRule правило периодичности
type FrequencyRule struct {
    Type           FrequencyType           `json:"type" validate:"required"`
    Daily          *DailyFrequency         `json:"daily,omitempty"`
    Monthly        *MonthlyFrequency       `json:"monthly,omitempty"`
    SpecificDates  *SpecificDatesFrequency `json:"specific_dates,omitempty"`
    Parity         *ParityFrequency        `json:"parity,omitempty"`
}

func (rule *FrequencyRule) IsEmpty() (bool) {
    return rule == nil || reflect.DeepEqual(*rule, FrequencyRule{})
}

func (rule *FrequencyRule) Valid() (bool, error) {
    
    if rule.IsEmpty() {
        return true, nil
    }

    switch rule.Type {

    case Daily:
        if rule.Daily == nil {
            return false, ErrDailyConfigIsRequired
        }
        if rule.Daily.Interval < 1 || rule.Daily.Interval > 365 {
            return false, ErrInvalidDailyInterval
        }
        if rule.Daily.Hour < 0 || rule.Daily.Hour > 23 {
            return false, ErrInvalidHour
        }
        if rule.Daily.Minute < 0 || rule.Daily.Minute > 59 {
            return false, ErrInvalidMinute
        }
        
    case Monthly:
        if rule.Monthly == nil {
            return false, ErrMonthlyConfigIsRequired
        }
        if rule.Monthly.DayOfMonth < 1 || rule.Monthly.DayOfMonth > 31 {
            return false, ErrInvalidMonthlyInterval
        }
        if rule.Monthly.Hour < 0 || rule.Monthly.Hour > 23 {
            return false, ErrInvalidHour
        }
        if rule.Monthly.Minute < 0 || rule.Monthly.Minute > 59 {
            return false, ErrInvalidMinute
        }
        
    case SpecificDates:
        if rule.SpecificDates == nil {
            return false, ErrSpecificDatesIsRequired
        }
        if len(rule.SpecificDates.Dates) == 0 {
            return false, ErrAtLeastOneDateIsRequired
        }
        
        for _, date := range rule.SpecificDates.Dates {
            // Парсим только дату (без времени)
            parsed, err := time.Parse("2006-01-02", date)
            if err != nil || parsed.Format("2006-01-02") != date {
                return false, &InvalidDateError{
                    Format: "YYYY-MM-DD",
                    Date:   date,
                    Err:    ErrInvalidDate,
                }
            }
        }
        
        // Дополнительная валидация часа и минуты
        if rule.SpecificDates.Hour < 0 || rule.SpecificDates.Hour > 23 {
            return false, ErrInvalidHour
        }
        if rule.SpecificDates.Minute < 0 || rule.SpecificDates.Minute > 59 {
            return false, ErrInvalidMinute
        }
    case Parity:
        if rule.Parity == nil {
            return false, ErrParityConfigIsRequired
        }
        if rule.Parity.Type != Even && rule.Parity.Type != Odd {
            return false, ErrInvalidParityType
        }
        
    default:
        return false, &UnknownFrequencyTypeError{
            Type: string(rule.Type),
            Err: ErrUnknownFrequencyType,
        }
    }
    
    return true, nil
}