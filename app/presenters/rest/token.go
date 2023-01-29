package rest

import (
	"net/http"
	"strings"

	"github.com/delveper/mylib/app/models"
)

type Token interface {
	models.AccessToken | models.RefreshToken
}

func retrieveToken[T Token](req *http.Request) (token T) {
	val := req.Context().Value(tokenContextKey)

	token, ok := val.(T)
	if ok {
		return token
	}

	return
}

func retrieveJWT(req *http.Request) string {
	for _, fn := range []func(*http.Request) string{
		tokenFromHeader,
		tokenFromURL,
		tokenFromCookie,
	} {
		if val := fn(req); val != "" {
			return val
		}
	}

	return ""
}

func tokenFromHeader(req *http.Request) string {
	header := req.Header.Get("Authorization")
	if len(header) > len(bearer) && strings.ToLower(header[:len(bearer)]) == bearer {
		return header[len(bearer)+1:]
	}

	return ""
}

func tokenFromCookie(req *http.Request) string {
	cookie, err := req.Cookie(refreshTokenKey)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func tokenFromURL(req *http.Request) string {
	return req.URL.Query().Get(refreshTokenKey)
}
