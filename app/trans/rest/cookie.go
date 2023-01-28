package rest

import (
	"net/http"
	"os"
	"time"
)

func setCookie(rw http.ResponseWriter, name, val string, exp time.Duration, age time.Duration) {
	http.SetCookie(rw, &http.Cookie{
		Name:     name,
		Value:    val,
		Domain:   os.Getenv("SRV_HOST"),
		Path:     "/auth",
		MaxAge:   int(age.Seconds()),
		Expires:  time.Now().Add(exp),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
}
