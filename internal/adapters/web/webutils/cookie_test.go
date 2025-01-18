package webutils

import (
	"net/http"
	"testing"
	"time"
)

func TestCookieStore(t *testing.T) {
	dataRetention := 24 * time.Hour
	store := CookieStore(dataRetention)

	if store == nil {
		t.Fatalf("expected non-nil store, got nil")
	}

	options := store.Options
	if options == nil {
		t.Fatalf("expected non-nil options, got nil")
	}

	if options.Path != "/" {
		t.Errorf("expected Path '/', got %s", options.Path)
	}

	if options.MaxAge != int(dataRetention.Seconds()) {
		t.Errorf("expected MaxAge %d, got %d", int(dataRetention.Seconds()), options.MaxAge)
	}

	if !options.HttpOnly {
		t.Errorf("expected HttpOnly true, got false")
	}

	if options.Secure {
		t.Errorf("expected Secure false, got true")
	}

	if options.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite %v, got %v", http.SameSiteLaxMode, options.SameSite)
	}
}
