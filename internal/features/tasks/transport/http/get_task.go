package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Get a task by ID
// @Description Get a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} GetTaskResponse "Task"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TaskHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	taskId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskId)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
