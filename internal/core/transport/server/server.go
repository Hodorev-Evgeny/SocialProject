package core_transport_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     ServerConfig
	log        *core_logger.Logger
	middleware []core_middleware.Middleware
}

func NewServer(
	config ServerConfig,
	log *core_logger.Logger,
	middleware ...core_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) ResisterApiVersionRouter(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.version)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}

}

func (s *HTTPServer) Start(ctx context.Context) error {
	mux := core_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error)

	go func() {
		defer close(ch)

		s.log.Warn("Starting HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server error: %w", err)
		}

	case <-ctx.Done():
		s.log.Warn("Stopping HTTP server", zap.Error(ctx.Err()))

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown server error: %w", err)
		}
		s.log.Warn("HTTP shutdown server stopped")
	}

	return nil
}
