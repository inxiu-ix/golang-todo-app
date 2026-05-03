package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_http_server "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/server"
)

type TaskHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		userId *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)
	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)
	DeleteTask(
		ctx context.Context,
		id int,
	) error
	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TaskHTTPHandler {
	return &TaskHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TaskHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
