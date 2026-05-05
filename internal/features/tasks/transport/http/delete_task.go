package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [delete]
func (h *TaskHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path")
		return
	}

	err = h.tasksService.DeleteTask(ctx, taskId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	responseHandler.NoContentResponse()
}
