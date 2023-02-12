package rest

import "time"

const bearer = "bearer"
const accessTokenKey = "access_token"
const refreshTokenKey = "refresh_token"
const xRequestID = "X-Request-ID"

type contextKey int

const (
	tokenContextKey contextKey = iota
	requestContextKey
)

const queryTimeout = 3 * time.Second
