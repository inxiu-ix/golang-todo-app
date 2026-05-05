package tasks_transport_http

import (
	"net/http"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Buy milk"`
	Description  *string `json:"description" validate:"omitempty,min=3,max=1000" example:"Buy milk at the store"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the given title, description and author user ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Create task request"
// @Success 201 {object} CreateTaskResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [post]
func (h *TaskHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
