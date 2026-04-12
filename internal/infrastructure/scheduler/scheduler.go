package scheduler

import (
    "context"
    "log"
    "sync"
    "time"
	taskdomain "example.com/taskservice/internal/domain/task"
)

// TaskScheduler реализует интерфейс Scheduler
type TaskScheduler struct {
    repo       taskdomain.Repository
    calculator taskdomain.NextRunCalculator
    executor   *TaskExecutor
    ticker     *time.Ticker
    stopChan   chan struct{}
    wg         sync.WaitGroup
    interval   time.Duration
}

type TaskExecutor struct {
    workerCount int
    taskQueue   chan taskdomain.Task
    wg          sync.WaitGroup
    repo        taskdomain.Repository
    calculator  taskdomain.NextRunCalculator
}

type Config struct {
    CheckInterval time.Duration
    WorkerCount   int
}

func NewTaskScheduler(repo taskdomain.Repository, config Config) *TaskScheduler {
    if config.CheckInterval == 0 {
        config.CheckInterval = 1 * time.Minute
    }
    if config.WorkerCount == 0 {
        config.WorkerCount = 5
    }

    return &TaskScheduler{
        repo:       repo,
        calculator: NewNextRunCalculator(),
        executor:   NewTaskExecutor(repo, NewNextRunCalculator(), config.WorkerCount),
        ticker:     time.NewTicker(config.CheckInterval),
        stopChan:   make(chan struct{}),
        interval:   config.CheckInterval,
    }
}

func NewTaskExecutor(repo taskdomain.Repository, calculator taskdomain.NextRunCalculator, workerCount int) *TaskExecutor {
    return &TaskExecutor{
        workerCount: workerCount,
        taskQueue:   make(chan taskdomain.Task, 100),
        repo:        repo,
        calculator:  calculator,
    }
}

func (s *TaskScheduler) Start(ctx context.Context) error {
    log.Println("Task scheduler starting...")

    if err := s.InitializeNextRunTimes(ctx); err != nil {
        log.Printf("Warning: failed to initialize next run times: %v", err)
    }

    s.executor.Start(ctx)
    
    s.wg.Add(1)
    go s.ScheduleLoop(ctx)

    return nil
}

func (s *TaskScheduler) Stop() {
    log.Println("Stopping task scheduler...")
    s.ticker.Stop()
    close(s.stopChan)
    s.wg.Wait()
    s.executor.Stop()
    log.Println("Task scheduler stopped")
}

func (s *TaskScheduler) ScheduleTask(ctx context.Context, task taskdomain.Task) (*taskdomain.Task, error) {
    nextRun, err := s.calculator.Calculate(task, time.Now())
    if err != nil {
        return nil, err
    }
    
    return s.repo.UpdateNextRunTime(ctx, task.ID, nextRun)
}

func (s *TaskScheduler) UnscheduleTask(ctx context.Context, taskID int64) error {
    return nil
}

func (s *TaskScheduler) InitializeNextRunTimes(ctx context.Context) error {
    tasks, err := s.repo.GetTasksForExecution(ctx)
    if err != nil {
        return err
    }

    for _, task := range tasks {
        if task.NextRunTime.IsZero() {
            nextRun, err := s.calculator.Calculate(task, time.Now())
            if err != nil {
                log.Printf("Failed to calculate next run for task %d: %v", task.ID, err)
                continue
            }
            if _, err := s.repo.UpdateNextRunTime(ctx, task.ID, nextRun); err != nil {
                log.Printf("Failed to update next run time for task %d: %v", task.ID, err)
            }
        }
    }

    return nil
}

func (s *TaskScheduler) ScheduleLoop(ctx context.Context) {
    defer s.wg.Done()

    for {
        select {
        case <-s.ticker.C:
            s.CheckAndScheduleTasks(ctx)
        case <-s.stopChan:
            return
        case <-ctx.Done():
            return
        }
    }
}

func (s *TaskScheduler) CheckAndScheduleTasks(ctx context.Context) {
    tasks, err := s.repo.GetTasksForExecution(ctx)
    if err != nil {
        log.Printf("Error getting tasks for execution: %v", err)
        return
    }

    for _, task := range tasks {
        select {
        case s.executor.taskQueue <- task:
        default:
            log.Printf("Task queue is full, dropping task %d", task.ID)
        }
    }
}

func (e *TaskExecutor) Start(ctx context.Context) {
    for i := 0; i < e.workerCount; i++ {
        e.wg.Add(1)
        go e.Worker(ctx, i)
    }
}

func (e *TaskExecutor) Worker(ctx context.Context, workerID int) {
    defer e.wg.Done()

    for {
        select {
        case task, ok := <-e.taskQueue:
            if !ok {
                return
            }
            e.ExecuteTask(ctx, task, workerID)
        case <-ctx.Done():
            return
        }
    }
}

func (e *TaskExecutor) ExecuteTask(ctx context.Context, task taskdomain.Task, workerID int) {
    log.Printf("[Worker %d] Executing task: %s (ID: %d, Status: %s)",
        workerID, task.Title, task.ID, task.Status)

    nextRun, err := e.calculator.Calculate(task, task.NextRunTime)
    if err != nil {
        log.Printf("[Worker %d] Failed to calculate next run time for task %d: %v", workerID, task.ID, err)
        return
    }
    
    updatedTask, err := e.repo.UpdateNextRunTime(ctx, task.ID, nextRun)
    if err != nil {
        log.Printf("[Worker %d] Failed to update next run time for task %d: %v", workerID, task.ID, err)
        return
    }
    
    log.Printf("[Worker %d] Task %s completed. Next run scheduled at: %v", 
        workerID, task.Title, updatedTask.NextRunTime)
}

func (e *TaskExecutor) Stop() {
    close(e.taskQueue)
    e.wg.Wait()
}