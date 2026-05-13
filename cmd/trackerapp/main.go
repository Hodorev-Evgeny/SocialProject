package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_pgx_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
	features_admin_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/admin/repository/postgres"
	feature_admin_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/admin/service"
	features_admin_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/admin/transport/http"
	features_users_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/repository/postgres"
	feature_user_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/service"
	features_users_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()
	fmt.Println("starting server")

	config := core_logger.MustNewConfig()
	logger, err := core_logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("starting expenses app")

	logger.Debug("starting initialization pool connection")
	pgconfig := core_pgx_pool.MustPostgresConfig()
	pool := core_pgx_pool.CreatePoolMust(ctx, pgconfig)
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("error pinging pool", zap.Error(err))
		os.Exit(1)
	}

	logger.Debug("starting initialization user service")
	userRepo := features_users_repository.NewUserRepository(pool)
	userServ := feature_user_service.NewUserService(userRepo)

	logger.Debug("starting initialization user transport")
	userTransporthttp := features_users_transport.NewUserHTTPHandler(userServ)
	userRouters := userTransporthttp.Routers()

	logger.Debug("starting initialization admin service")
	adminRepo := features_admin_repository.NewAdminRepository(pool)
	adminServ := feature_admin_service.NewAdminService(adminRepo)

	logger.Debug("starting initialization admin transport")
	adminTransporthttp := features_admin_transport.NewAdminHTTPHandler(adminServ)
	adminRouters := adminTransporthttp.Routers()

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRouters...)
	apiVersionRouter.RegisterRoutes(adminRouters...)

	httpServer := core_transport_server.NewServer(
		core_transport_server.MustNewConfigServer(),
		logger,
		core_middleware.RequestId(),
		core_middleware.Logger(logger),
		core_middleware.Trace(),
		core_middleware.PanicRecovery(),
	)
	httpServer.ResisterApiVersionRouter(apiVersionRouter)

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server failed to start", zap.Error(err))
	}
}
