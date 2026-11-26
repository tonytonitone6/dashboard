package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"dashboard/api/internal/auth"
	"dashboard/api/internal/server"
)

func TestTheAssembledAPIIsBehindTheGate(t *testing.T) {
	env := map[string]string{
		"PUBLIC_URL":           "https://dashboard.example",
		"SESSION_KEY":          strings.Repeat("k", 32),
		"GOOGLE_CLIENT_ID":     "client-id",
		"GOOGLE_CLIENT_SECRET": "client-secret",
		"OWNER_GOOGLE_SUB":     "1001",
	}
	cfg, err := auth.ConfigFromEnv(func(name string) string { return env[name] })
	if err != nil {
		t.Fatalf("test configuration: %v", err)
	}
	api := httptest.NewServer(server.New(cfg, time.Now))
	t.Cleanup(api.Close)

	for _, tc := range []struct {
		path string
		want int
	}{
		{"/api/health", http.StatusOK},
		{"/api/v1/overview", http.StatusUnauthorized},
		{"/api/v1/me", http.StatusUnauthorized},
		{"/api/v1/does-not-exist", http.StatusUnauthorized},
	} {
		resp, err := http.Get(api.URL + tc.path)
		if err != nil {
			t.Fatalf("GET %s: %v", tc.path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != tc.want {
			t.Errorf("GET %s without a session = %d, want %d", tc.path, resp.StatusCode, tc.want)
		}
	}
}
