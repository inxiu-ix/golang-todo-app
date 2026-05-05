package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Get tasks
// @Description Get tasks with optional user ID, limit and offset
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_id query int false "User ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} GetTasksResponse "Tasks"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, limit, offset, err := getUserIDLimitOffsetOueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user ID, limit and offset query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userId, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetOueryParams(
	r *http.Request,
) (*int, *int, *int, error) {
	const (
		userIdQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	userId, err := core_http_request.GetIntQueryParam(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user ID query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	return userId, limit, offset, nil
}
