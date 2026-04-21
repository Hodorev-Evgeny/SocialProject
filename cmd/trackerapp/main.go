package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
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

	userTransporthttp := features_users_transport.NewUserHTTPHandler(nil)
	userRouters := userTransporthttp.Routers()

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(userRouters...)

	httpServer := core_transport_server.NewServer(
		core_transport_server.MustNewConfigServer(),
		logger,
		core_middleware.RequestId(),
		core_middleware.Logger(logger),
		core_middleware.PanicRecovery(),
		core_middleware.Trace(),
	)
	httpServer.ResisterApiVersionRouter(apiVersionRouter)

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server failed to start", zap.Error(err))
	}
}
