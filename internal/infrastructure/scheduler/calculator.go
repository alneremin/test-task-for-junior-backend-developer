package scheduler

import (
    "fmt"
	"time"
    taskdomain "example.com/taskservice/internal/domain/task"
)

type NextRunCalculatorImpl struct{}

func NewNextRunCalculator() *NextRunCalculatorImpl {
    return &NextRunCalculatorImpl{}
}

// Calculate реализует taskdomain.NextRunCalculator интерфейс
func (c *NextRunCalculatorImpl) Calculate(task taskdomain.Task, now time.Time) (time.Time, error) {

    switch task.Frequency.Type {
    case taskdomain.Daily:
        return c.CalculateDaily(*task.Frequency.Daily, now)
    case taskdomain.Monthly:
        return c.CalculateMonthly(*task.Frequency.Monthly, now)
    case taskdomain.SpecificDates:
        return c.CalculateSpecificDates(*task.Frequency.SpecificDates, now)
    case taskdomain.Parity:
        return c.CalculateParity(*task.Frequency.Parity, now)
    default:
        return time.Time{}, taskdomain.ErrUnknownFrequencyType
    }
}


func (c *NextRunCalculatorImpl) CalculateDaily(freq taskdomain.DailyFrequency, fromTime time.Time) (time.Time, error) {
    // Устанавливаем время на указанный час и минуту для текущего дня
    nextRun := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(),
        freq.Hour, freq.Minute, 0, 0, fromTime.Location())
    
    // Если заданное время уже прошло сегодня, начинаем со следующего дня
    if !nextRun.After(fromTime) {
        nextRun = nextRun.AddDate(0, 0, 1)
    }
    
    // Прибавляем интервал в днях
    // Если интервал = 1, то прибавляем 1 день
    // Если интервал = 3, то прибавляем 3 дня и т.д.
    if freq.Interval > 1 {
        nextRun = nextRun.AddDate(0, 0, freq.Interval-1)
    }
    
    return nextRun, nil
}

func (c *NextRunCalculatorImpl) CalculateMonthly(freq taskdomain.MonthlyFrequency, fromTime time.Time) (time.Time, error) {
    // Получаем первый день следующего месяца, содержащего нужное число
    year := fromTime.Year()
    month := fromTime.Month()
    
    // Пытаемся найти дату в текущем месяце
    targetDay := freq.DayOfMonth
    
    // Проверяем, существует ли такой день в текущем месяце
    lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, fromTime.Location()).Day()
    
    // Корректируем день, если он выходит за пределы месяца
    if targetDay > lastDayOfMonth {
        targetDay = lastDayOfMonth
    }
    
    nextRun := time.Date(year, month, targetDay,
        freq.Hour, freq.Minute, 0, 0, fromTime.Location())
    
    // Если время в текущем месяце уже прошло, переходим на следующий месяц
    if nextRun.Before(fromTime) {
        // Переходим на следующий месяц
        nextMonth := fromTime.AddDate(0, 1, 0)
        year = nextMonth.Year()
        month = nextMonth.Month()
        
        // Повторно проверяем существование дня в следующем месяце
        lastDayOfMonth = time.Date(year, month+1, 0, 0, 0, 0, 0, fromTime.Location()).Day()
        if targetDay > lastDayOfMonth {
            targetDay = lastDayOfMonth
        }
        
        nextRun = time.Date(year, month, targetDay,
            freq.Hour, freq.Minute, 0, 0, fromTime.Location())
    }
    
    return nextRun, nil
}

func (c *NextRunCalculatorImpl) CalculateSpecificDates(freq taskdomain.SpecificDatesFrequency, fromTime time.Time) (time.Time, error) {
    // Парсим все даты и сортируем их
    var dates []time.Time
    
    for _, dateStr := range freq.Dates {
        parsedDate, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", dateStr)
        }
        
        // Устанавливаем время выполнения
        dateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
            freq.Hour, freq.Minute, 0, 0, fromTime.Location())
        
        // Добавляем только будущие даты
        if dateTime.After(fromTime) {
            dates = append(dates, dateTime)
        }
    }
    
    if len(dates) == 0 {
        return time.Time{}, fmt.Errorf("no future dates available for execution")
    }
    
    // Находим ближайшую дату
    nextRun := dates[0]
    for _, date := range dates {
        if date.Before(nextRun) {
            nextRun = date
        }
    }
    
    return nextRun, nil
}

func (c *NextRunCalculatorImpl) CalculateParity(freq taskdomain.ParityFrequency, fromTime time.Time) (time.Time, error) {
    // Начинаем с текущей даты
    nextRun := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(),
        freq.Hour, freq.Minute, 0, 0, fromTime.Location())
    
    // Если время уже прошло сегодня, начинаем с завтра
    if nextRun.Before(fromTime) || nextRun.Equal(fromTime) {
        nextRun = nextRun.AddDate(0, 0, 1)
    }
    
    // Ищем подходящий день с нужной четностью
    for {
        day := nextRun.Day()
        isEven := day%2 == 0
        
        // Проверяем соответствие четности
        if (freq.Type == taskdomain.Even && isEven) || (freq.Type == taskdomain.Odd && !isEven) {
            break
        }
        
        // Переходим к следующему дню
        nextRun = nextRun.AddDate(0, 0, 1)
    }
    
    return nextRun, nil
}
