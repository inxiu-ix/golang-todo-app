package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, from, to, err := getStatisticsQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}
	response := GetStatisticsResponse(toDTOFromDomain(statistics))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistic) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getStatisticsQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdQueryParamKey   = "user_id"
		fromDateQueryParamKey = "from"
		toDateQueryParamKey   = "to"
	)

	userId, err := core_http_request.GetIntQueryParam(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user ID query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromDateQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get from date query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toDateQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get to date query param: %w", err)
	}

	return userId, from, to, nil

}
