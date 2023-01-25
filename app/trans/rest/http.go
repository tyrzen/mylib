package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/delveper/mylib/app/ent"
)

type response struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandlerLoggerFunc main func that will be used for all handlers.
type HandlerLoggerFunc func(http.ResponseWriter, *http.Request, ent.Logger)

// ServeHTTP gives handlerLoggerFunc feature of http.Handler.
// ps. don't be dogmatic about injecting logger into context.
func (hlf HandlerLoggerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hlf(rw, req, loggerFromContext(rw, req))
}

func loggerFromContext(rw http.ResponseWriter, req *http.Request) ent.Logger {
	logger, ok := req.Context().Value(loggerContextKey).(ent.Logger)
	if !ok {
		respond(rw, req, http.StatusInternalServerError, ErrLoggerNotFound)
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
func setCookie(rw http.ResponseWriter, name, val string, exp time.Duration, age int) {
	http.SetCookie(rw, &http.Cookie{
		Name:     name,
		Value:    val,
		Domain:   os.Getenv("SRV_HOST"),
		Path:     "/auth",
		MaxAge:   age,
		Expires:  time.Now().Add(exp),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
}

func respond(rw http.ResponseWriter, req *http.Request, code int, data any) {
	logger := loggerFromContext(rw, req)

	if data == nil && code != http.StatusNoContent {
		logger.Errorw("Failed writing response due nil data.",
			"object", nil,
			"error", ErrInvalidData,
		)
		respond(rw, req, http.StatusBadRequest, ErrInvalidData)

		return
	}

	if err, ok := data.(error); ok {
		data = response{Message: http.StatusText(code), Details: err.Error()}
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		logger.Errorw("Failed encoding data to JSON.",
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

func tokenFromRequest(req *http.Request) (token ent.Token) {
	for _, fn := range []func(*http.Request) string{tokenFromHeader, tokenFromCookie, tokenFromContext} {
		if val := fn(req); val != "" {
			return ent.Token{Value: val}
		}
	}

	return
}

func tokenFromHeader(req *http.Request) string {
	header := req.Header.Get("Authorization")
	if len(header) > len(bearer) && strings.ToLower(header[:len(bearer)]) == bearer {
		return header[len(bearer):]
	}

	return ""
}

func tokenFromCookie(req *http.Request) string {
	cookie, err := req.Cookie(tokenCookieKey)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func tokenFromContext(req *http.Request) string {
	ctx := req.Context()
	val := ctx.Value(tokenContextKey)
	token, ok := val.(string)
	if !ok || val == "" {

		return ""
	}

	return token
}
