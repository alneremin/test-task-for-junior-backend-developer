package task

import (
    "testing"
)

func TestFrequencyRule_Valid(t *testing.T) {
    tests := []struct {
        name    string
        rule    *FrequencyRule
        want    bool
        wantErr bool
    }{
        {
            name: "valid daily frequency",
            rule: &FrequencyRule{
                Type: Daily,
                Daily: &DailyFrequency{
                    Interval: 5,
                    Hour:     10,
                    Minute:   30,
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "valid monthly frequency",
            rule: &FrequencyRule{
                Type: Monthly,
                Monthly: &MonthlyFrequency{
                    DayOfMonth: 15,
                    Hour:       14,
                    Minute:     0,
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "valid specific dates frequency",
            rule: &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates: []string{"2024-12-25", "2024-12-31"},
                    Hour:  12,
                    Minute: 0,
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "valid parity even",
            rule: &FrequencyRule{
                Type: Parity,
                Parity: &ParityFrequency{
                    Type:   Even,
                    Hour:   8,
                    Minute: 0,
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "valid parity odd",
            rule: &FrequencyRule{
                Type: Parity,
                Parity: &ParityFrequency{
                    Type:   Odd,
                    Hour:   17,
                    Minute: 45,
                },
            },
            want:    true,
            wantErr: false,
        },
        {
            name: "daily config missing",
            rule: &FrequencyRule{
                Type:  Daily,
                Daily: nil,
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "daily interval too small",
            rule: &FrequencyRule{
                Type: Daily,
                Daily: &DailyFrequency{
                    Interval: 0,
                    Hour:     10,
                    Minute:   30,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "daily interval too large",
            rule: &FrequencyRule{
                Type: Daily,
                Daily: &DailyFrequency{
                    Interval: 400,
                    Hour:     10,
                    Minute:   30,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "monthly config missing",
            rule: &FrequencyRule{
                Type:    Monthly,
                Monthly: nil,
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "monthly day too small",
            rule: &FrequencyRule{
                Type: Monthly,
                Monthly: &MonthlyFrequency{
                    DayOfMonth: 0,
                    Hour:       14,
                    Minute:     0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "monthly day too large",
            rule: &FrequencyRule{
                Type: Monthly,
                Monthly: &MonthlyFrequency{
                    DayOfMonth: 32,
                    Hour:       14,
                    Minute:     0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "specific dates config missing",
            rule: &FrequencyRule{
                Type:          SpecificDates,
                SpecificDates: nil,
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "specific dates empty list",
            rule: &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates: []string{},
                    Hour:  12,
                    Minute: 0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "specific dates invalid format",
            rule: &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates: []string{"2024-12-25 12:00"},
                    Hour:  12,
                    Minute: 0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "specific dates completely wrong format",
            rule: &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates: []string{"invalid-date"},
                    Hour:  12,
                    Minute: 0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "parity config missing",
            rule: &FrequencyRule{
                Type:   Parity,
                Parity: nil,
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "parity invalid type",
            rule: &FrequencyRule{
                Type: Parity,
                Parity: &ParityFrequency{
                    Type:   "invalid",
                    Hour:   8,
                    Minute: 0,
                },
            },
            want:    false,
            wantErr: true,
        },
        {
            name: "unknown frequency type",
            rule: &FrequencyRule{
                Type: "unknown",
            },
            want:    false,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.rule.Valid()
            if (err != nil) != tt.wantErr {
                t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Valid() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestDailyFrequency_EdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        interval int
        hour     int
        minute   int
        wantErr  bool
    }{
        {"min valid interval", 1, 0, 0, false},
        {"max valid interval", 365, 23, 59, false},
        {"interval 366", 366, 12, 30, true},
        {"hour -1", 5, -1, 30, true}, // hour validation not in Valid()
        {"hour 24", 5, 24, 30, true}, // hour validation not in Valid()
        {"minute -1", 5, 12, -1, true},
        {"minute 60", 5, 12, 60, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            rule := &FrequencyRule{
                Type: Daily,
                Daily: &DailyFrequency{
                    Interval: tt.interval,
                    Hour:     tt.hour,
                    Minute:   tt.minute,
                },
            }
            _, err := rule.Valid()
            if (err != nil) != tt.wantErr {
                t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestSpecificDatesFrequency_DateFormats(t *testing.T) {
    validDates := []string{
        "2024-01-01",
        "2024-12-31",
        "2024-02-29", // leap year
        "2023-02-28",
        "2025-06-15",
    }

    for _, date := range validDates {
        t.Run("valid date: "+date, func(t *testing.T) {
            rule := &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates:  []string{date},
                    Hour:   12,
                    Minute: 0,
                },
            }
            got, err := rule.Valid()
            if err != nil {
                t.Errorf("Expected no error for valid date %s, got %v", date, err)
            }
            if !got {
                t.Errorf("Expected true for valid date %s, got false", date)
            }
        })
    }

    invalidDates := []string{
        "2024-13-01", // invalid month
        "2024-00-01", // invalid month
        "2024-01-32", // invalid day
        "2024-02-30", // invalid day for feb
        "not-a-date",
    }

    for _, date := range invalidDates {
        t.Run("invalid date: "+date, func(t *testing.T) {
            rule := &FrequencyRule{
                Type: SpecificDates,
                SpecificDates: &SpecificDatesFrequency{
                    Dates:  []string{date},
                    Hour:   12,
                    Minute: 0,
                },
            }
            _, err := rule.Valid()
            if err == nil {
                t.Errorf("Expected error for invalid date %s, got nil", date)
            }
        })
    }
}

func TestParityFrequency_Validation(t *testing.T) {
    validParityTypes := []ParityType{Even, Odd}
    
    for _, parityType := range validParityTypes {
        t.Run("valid parity type: "+string(parityType), func(t *testing.T) {
            rule := &FrequencyRule{
                Type: Parity,
                Parity: &ParityFrequency{
                    Type:   parityType,
                    Hour:   12,
                    Minute: 0,
                },
            }
            got, err := rule.Valid()
            if err != nil {
                t.Errorf("Unexpected error for valid parity type: %v", err)
            }
            if !got {
                t.Errorf("Expected true for valid parity type, got false")
            }
        })
    }

    t.Run("invalid parity type", func(t *testing.T) {
        rule := &FrequencyRule{
            Type: Parity,
            Parity: &ParityFrequency{
                Type:   "invalid",
                Hour:   12,
                Minute: 0,
            },
        }
        _, err := rule.Valid()
        if err == nil {
            t.Error("Expected error for invalid parity type, got nil")
        }
    })
}

func TestFrequencyRule_MultipleDates(t *testing.T) {
    t.Run("multiple valid dates", func(t *testing.T) {
        rule := &FrequencyRule{
            Type: SpecificDates,
            SpecificDates: &SpecificDatesFrequency{
                Dates: []string{
                    "2024-01-01",
                    "2024-02-01",
                    "2024-03-01",
                    "2024-04-01",
                },
                Hour:   9,
                Minute: 0,
            },
        }
        got, err := rule.Valid()
        if err != nil {
            t.Errorf("Unexpected error: %v", err)
        }
        if !got {
            t.Error("Expected true for multiple valid dates")
        }
    })

    t.Run("mixed valid and invalid dates", func(t *testing.T) {
        rule := &FrequencyRule{
            Type: SpecificDates,
            SpecificDates: &SpecificDatesFrequency{
                Dates: []string{
                    "2024-01-01 09:00",
                    "invalid-date",
                    "2024-03-01 09:00",
                },
                Hour:   9,
                Minute: 0,
            },
        }
        _, err := rule.Valid()
        if err == nil {
            t.Error("Expected error for mixed valid/invalid dates")
        }
    })
}

func TestFrequencyRule_NilConfigs(t *testing.T) {
    frequencyTypes := []FrequencyType{Daily, Monthly, SpecificDates, Parity}
    
    for _, freqType := range frequencyTypes {
        t.Run("nil config for "+string(freqType), func(t *testing.T) {
            rule := &FrequencyRule{
                Type: freqType,
            }
            
            // Set all configs to nil
            rule.Daily = nil
            rule.Monthly = nil
            rule.SpecificDates = nil
            rule.Parity = nil
            
            _, err := rule.Valid()
            if err == nil {
                t.Errorf("Expected error for nil config with type %s", freqType)
            }
        })
    }
}

// Benchmark tests
func BenchmarkFrequencyRule_Valid(b *testing.B) {
    rule := &FrequencyRule{
        Type: Daily,
        Daily: &DailyFrequency{
            Interval: 5,
            Hour:     10,
            Minute:   30,
        },
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rule.Valid()
    }
}

func BenchmarkSpecificDates_Valid(b *testing.B) {
    dates := make([]string, 100)
    for i := 0; i < 100; i++ {
        dates[i] = "2024-01-01 10:00"
    }
    
    rule := &FrequencyRule{
        Type: SpecificDates,
        SpecificDates: &SpecificDatesFrequency{
            Dates:  dates,
            Hour:   10,
            Minute: 0,
        },
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rule.Valid()
    }
}