package rest

import (
	"net/http"
	"strings"

	"github.com/delveper/mylib/app/ent"
)

func retrieveToken(req *http.Request) (token ent.Token) {
	for _, fn := range []func(*http.Request) string{
		tokenFromHeader,
		tokenFromCookie,
		tokenFromContext,
		tokenFromURL,
	} {
		if val := fn(req); val != "" {
			return ent.Token{Value: val}
		}
	}

	return
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

func tokenFromContext(req *http.Request) string {
	val := req.Context().Value(tokenContextKey)
	token, ok := val.(string)
	if ok && val != "" {
		return token
	}

	return ""
}

func tokenFromURL(req *http.Request) string {
	return req.URL.Query().Get(refreshTokenKey)
}
