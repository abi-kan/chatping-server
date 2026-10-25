package core_http_response

import (
	core_errors "chatping-server/internal/core/errors"
	core_logger "chatping-server/internal/core/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ResponseHandler struct {
	logger         *core_logger.Logger
	responseWriter http.ResponseWriter
}

func NewResponseHandler(
	logger *core_logger.Logger,
	responseWriter http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		logger:         logger,
		responseWriter: responseWriter,
	}
}

func (rh *ResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	rh.responseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(rh.responseWriter).Encode(responseBody); err != nil {
		rh.logger.Error("Write HTTP response", zap.Error(err))
	}
}

func (rh *ResponseHandler) NoContentResponse() {
	rh.responseWriter.WriteHeader(http.StatusNoContent)
}

func (rh *ResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = rh.logger.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = rh.logger.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = rh.logger.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = rh.logger.Error
	}

	logFunc(msg, zap.Error(err))
	rh.errorResponse(statusCode, err, msg)
}

func (rh *ResponseHandler) PanicResponse(anyPanic any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", anyPanic)

	rh.logger.Error(msg, zap.Error(err))
	rh.errorResponse(statusCode, err, msg)
}

func (rh *ResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	rh.JSONResponse(response, statusCode)
}
