package webutils

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	keyLength = 16
)

func CookieStore(dataRetention time.Duration) *sessions.CookieStore {
	cookieStore := sessions.NewCookieStore(securecookie.GenerateRandomKey(keyLength), securecookie.GenerateRandomKey(keyLength))

	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(dataRetention.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return cookieStore
}
