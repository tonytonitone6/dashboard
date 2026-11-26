// Package server assembles the API: every route, behind the gate.
package server

import (
	"net/http"
	"time"

	"dashboard/api/internal/auth"
	"dashboard/api/internal/overview"
)

// New is the only place routes are registered, so the tests exercise the same
// handler that production serves.
func New(cfg auth.Config, now func() time.Time) http.Handler {
	sessions := auth.New(cfg, now)

	mux := http.NewServeMux()
	mux.HandleFunc(auth.HealthPattern, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	sessions.Register(mux)
	overview.Register(mux)

	return sessions.Gate(mux)
}
