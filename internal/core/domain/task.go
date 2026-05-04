package domain

import (
	"fmt"
	"time"

	core_errors "github.com/inxiu-ix/golang-todo-app/internal/core/errors"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	AuthorUserID int
	CreatedAt    time.Time
	CompletedAt  *time.Time
	Completed    bool
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	authorUserID int,
	createdAt time.Time,
	completedAt *time.Time,
	completed bool,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		AuthorUserID: authorUserID,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		Completed:    completed,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UnitializedID,
		UnitializedVersion,
		title,
		description,
		authorUserID,
		time.Now(),
		nil,
		false,
	)
}

func (t *Task) CompletionDuration() *time.Duration {
	if !t.Completed {
		return nil
	}

	if t.CompletedAt == nil {
		return nil
	}
 
	duration := t.CompletedAt.Sub(t.CreatedAt)
	return &duration
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))

	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"title must be between 1 and 100 characters: %d: %w",
			titleLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))

		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf(
				"description must be between 1 and 1000 characters: %d: %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"completed_at must be set when task is completed: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"completed_at must be after created_at: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"completed_at must be null when task is not completed: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf(
			"`Completed` can't be patched to null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("failed to validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value

		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate task after patch: %w", err)
	}

	*t = tmp

	return nil
}
