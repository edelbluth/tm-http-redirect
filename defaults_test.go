package tm_http_redirect_test

import (
	"testing"

	"github.com/edelbluth/tm_http_redirect"
)

// This is only to ensure the constant values do not get changed accidently
func TestDefaultRedirectionStatusCode(t *testing.T) {
	if tm_http_redirect.DefaultRedirectionStatusCode != 307 {
		t.Fatalf("unexpected DefaultRedirectionStatusCode: %v", tm_http_redirect.DefaultRedirectionStatusCode)
	}
}

func TestDefaultRedirectionHeader(t *testing.T) {
	if tm_http_redirect.DefaultRedirectionHeader != "Location" {
		t.Fatalf("unexpected DefaultRedirectionHeader: %v", tm_http_redirect.DefaultRedirectionHeader)
	}
}
