package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

type Message struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandlerLoggerFunc main func that will be used for all handlers.
type HandlerLoggerFunc func(http.ResponseWriter, *http.Request, ent.Logger)

type LoggerKey string

const loggerKey LoggerKey = "logger" // var loggerKey = &struct{ string }{"logger"}

// ServeHTTP gives HandlerLoggerFunc http.Handler features.
// ps. don't be dogmatic about injecting logger into context.
func (hlf HandlerLoggerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger, ok := req.Context().Value(loggerKey).(ent.Logger)
	if !ok {
		http.Error(rw, "failed inject logger", http.StatusInternalServerError)
	}

	hlf(rw, req, logger)
}

// WithLogRequest logs every request
// and injects logger into context scope of request.
func WithLogRequest(logger ent.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"method", req.Method,
				"uri", req.RequestURI,
				"user-agent", req.UserAgent(),
				"remote", req.RemoteAddr,
			)

			ctx := context.WithValue(req.Context(), loggerKey, logger)

			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}

/*
func (resp *Responder) Write(statusCode int, data any) {
	if data == nil {
		if statusCode != http.StatusNoContent {
			resp.Errorf("Invalid data, expected nil")
		}

		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		resp.Errorw("Failed encoding to JSON.",
			"object", data,
			"error", err,
		)
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)

	if _, err := buf.WriteTo(resp.responseWriter); err != nil {
		resp.logger.Errorw("Failed writing response from buffer.",
			"object", data,
			"error", err,
		)
	}
}

func (resp *Responder) decodeBody(data any) error {
	defer func() {
		if err := resp.request.Body.Close(); err != nil {
			resp.logger.Warnf("error closing request body")
		}
	}()

	if err := json.NewDecoder(resp.request.Body).Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON body: %w", err)
	}

	return nil
}
*/
