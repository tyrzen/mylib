package rest

import "time"

const bearer = "bearer"

const xRequestID = "X-Request-ID"

const refreshTokenKey = "refresh_token"

const queryTimeout = 5 * time.Second

type contextKey int

const (
	loggerContextKey contextKey = iota
	tokenContextKey
	requestContextKey
	readerContextKey
)
