package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_pgx_pool "github.com/inxiu-ix/golang-todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/inxiu-ix/golang-todo-app/internal/features/tasks/repository/postgres"
	tasks_service "github.com/inxiu-ix/golang-todo-app/internal/features/tasks/service"
	tasks_transport_http "github.com/inxiu-ix/golang-todo-app/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/inxiu-ix/golang-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/inxiu-ix/golang-todo-app/internal/features/users/service"
	users_transport_http "github.com/inxiu-ix/golang-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres pool...")

	pool, err := core_pgx_pool.NewPool(core_pgx_pool.NewConfigMust(), ctx)
	if err != nil {
		logger.Fatal("failed to create postgres pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("time zone app", zap.Any("time zone", timeZone))

	logger.Debug("Initializing feature...", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	logger.Debug("initializing feature...", zap.String("feature", "tasks"))

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server...")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersionV1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
	}
}
