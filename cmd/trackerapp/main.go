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
	features_auth_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/repository/postgres"
	features_auth_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/service"
	features_auth_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/auth/transport/http"
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

	logger.Debug("starting initialization auth service")
	authRepo := features_auth_repository.NewAuthRepository(pool)
	jwtCfg := features_auth_service.MustJWTConfig()
	authServ := features_auth_service.NewAuthService(authRepo, userRepo, jwtCfg)

	logger.Debug("starting initialization user transport")
	userTransporthttp := features_users_transport.NewUserHTTPHandler(userServ, authServ)
	userRouters := userTransporthttp.Routers()

	logger.Debug("starting initialization auth transport")
	authTransporthttp := features_auth_transport.NewAuthHTTPHandler(authServ)
	authRouters := authTransporthttp.Routers()

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(append(userRouters, authRouters...)...)

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
