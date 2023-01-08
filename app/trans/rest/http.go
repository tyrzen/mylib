package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/pkg/errors"
)

type Response struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandlerLoggerFunc main func that will be used for all handlers.
type HandlerLoggerFunc func(http.ResponseWriter, *http.Request, ent.Logger)

type logKey int

const loggerKey logKey = iota // var loggerKey = &struct{}{}

// ServeHTTP gives HandlerLoggerFunc feature of http.Handler.
// ps. don't be dogmatic about injecting logger into context.
func (hlf HandlerLoggerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hlf(rw, req, extractLogger(rw, req))
}

func extractLogger(rw http.ResponseWriter, req *http.Request) ent.Logger {
	logger, ok := req.Context().Value(loggerKey).(ent.Logger)
	if !ok {
		respond(rw, req, http.StatusInternalServerError, errors.New("failed extracting logger from request"))
	}

	return logger
}

func decodeBody(req *http.Request, data any) (err error) {
	defer func() {
		if e := req.Body.Close(); e != nil {
			err = fmt.Errorf("error while closing request body: %w", err)
		}
	}()

	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON body: %w", err)
	}

	return nil
}

//nolint:godox    // TODO: Cancellation of database operations in case of respond errors.
func respond(rw http.ResponseWriter, req *http.Request, code int, data any) {
	logger := extractLogger(rw, req)

	if data == nil && code != http.StatusNoContent {
		logger.Errorw("Failed writing response due nil data.",
			"object", nil,
			"error", ErrInvalidData,
		)
		respond(rw, req, http.StatusBadRequest, ErrInvalidData)

		return
	}

	if err, ok := data.(error); ok {
		data = Response{Message: http.StatusText(code), Details: err.Error()}
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		logger.Errorw("Failed encoding to JSON.",
			"object", data,
			"error", err)
		respond(rw, req, http.StatusInternalServerError, ErrEncoding)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		logger.Errorw("Failed writing response from buffer.",
			"object", data,
			"error", fmt.Errorf("%w: %v", ErrWritingResponse, err),
		)
	}
}
