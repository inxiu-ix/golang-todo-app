package domain

import "time"

type Statistic struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistic(
	TasksCreated int,
	TasksCompleted int,
	TasksCompletedRate *float64,
	TasksAverageCompletionTime *time.Duration,
) Statistic {
	return Statistic{
		TasksCreated:               TasksCreated,
		TasksCompleted:             TasksCompleted,
		TasksCompletedRate:         TasksCompletedRate,
		TasksAverageCompletionTime: TasksAverageCompletionTime,
	}
}
