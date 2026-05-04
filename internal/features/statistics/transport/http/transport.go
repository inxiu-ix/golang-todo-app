package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_http_server "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistic, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method: http.MethodGet,
			Path:   "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
