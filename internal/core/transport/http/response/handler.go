package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"go.uber.org/zap"
)

type HandlerResponse struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHandlerResponse(logger *core_logger.Logger, rw http.ResponseWriter) *HandlerResponse {
	return &HandlerResponse{
		log: logger,
		rw:  rw,
	}
}

func (h *HandlerResponse) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logfunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrorBadRequest):
		statusCode = http.StatusBadRequest
		logfunc = h.log.Warn

	case errors.Is(err, core_errors.ErrorUnauthorized):
		statusCode = http.StatusUnauthorized
		logfunc = h.log.Warn

	case errors.Is(err, core_errors.ErrorValidation):
		statusCode = http.StatusUnprocessableEntity
		logfunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logfunc = h.log.Error
	}

	logfunc(msg, zap.Error(err))

	h.errorResponse(err, msg, statusCode)
}

func (h *HandlerResponse) PanicResponse(p any, msg string) {
	statuscode := http.StatusInternalServerError
	err := fmt.Errorf("unexepted punic:", p)

	h.log.Error(msg, zap.Error(err))
	h.errorResponse(err, msg, statuscode)
}

func (h *HandlerResponse) errorResponse(err error, msg string, status int) {
	h.rw.WriteHeader(status)

	response := map[string]string{
		"massage": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("write HTTP error response", zap.Error(err))
	}
}
