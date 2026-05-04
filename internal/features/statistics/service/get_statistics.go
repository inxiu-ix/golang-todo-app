package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_errors "github.com/inxiu-ix/golang-todo-app/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistic, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistic{}, fmt.Errorf("to date must be after from date: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return domain.Statistic{}, fmt.Errorf("failed to get tasks: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistic {
	if len(tasks) == 0 {
		return domain.Statistic{
			TasksCreated:               0,
			TasksCompleted:             0,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	totalTasks := len(tasks)
	completedTasks := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			completedTasks++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}
	completedRate := float64(completedTasks) / float64(totalTasks) * 100

	var tasksAverageCompletionTime *time.Duration
	if completedTasks > 0 && totalCompletionDuration != 0 {
		averageCompletionTime := totalCompletionDuration / time.Duration(completedTasks)
		tasksAverageCompletionTime = &averageCompletionTime
	}

	return domain.NewStatistic(
		totalTasks,
		completedTasks,
		&completedRate,
		tasksAverageCompletionTime,
	)

}
