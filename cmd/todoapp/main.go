package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_pgx_pool "github.com/inxiu-ix/golang-todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/middleware"
	core_http_server "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/inxiu-ix/golang-todo-app/internal/features/users/repository/postgres"
	users_service "github.com/inxiu-ix/golang-todo-app/internal/features/users/service"
	users_transport_http "github.com/inxiu-ix/golang-todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
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

	logger.Debug("Initializing features...", zap.String("features", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

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
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("failed to run HTTP server", zap.Error(err))
	}
}
