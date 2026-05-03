package tasks_service

import (
	"context"
	"fmt"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("failed to validate task: %w", err)
	}

	task, err := s.tasksRepository.CreateTask(ctx, task)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
