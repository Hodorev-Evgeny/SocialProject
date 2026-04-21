package response

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *HandlerResponse) PanicResponse(p any, msg string) {
	statuscode := http.StatusInternalServerError
	err := fmt.Errorf("unexepted punic:", p)

	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statuscode)

	response := map[string]string{
		"massage": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}
