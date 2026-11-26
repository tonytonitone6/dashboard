// Package auth locks the API to its owner: a deny-by-default gate, OAuth
// sign-in, and a stateless signed session cookie.
package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Identity is who a session belongs to. Provider and Subject are the permanent
// key. Name is display text and never takes part in an access decision.
type Identity struct {
	Provider string
	Subject  string
	Name     string
}

type Service struct {
	cfg    Config
	now    func() time.Time
	client *http.Client
}

func New(cfg Config, now func() time.Time) *Service {
	return &Service{
		cfg: cfg,
		now: now,
		client: &http.Client{
			Timeout: 10 * time.Second,
			// A redirect could carry the client secret, the code or the token to another URL.
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

// HealthPattern is the one public route this package does not serve itself.
const HealthPattern = "GET /api/health"

const (
	loginPattern    = "GET /api/auth/{provider}/login"
	callbackPattern = "GET /api/auth/{provider}/callback"
	logoutPattern   = "POST /api/auth/logout"
)

// publicPatterns is the whole list of routes that work without a session.
var publicPatterns = map[string]bool{
	HealthPattern:   true,
	loginPattern:    true,
	callbackPattern: true,
	logoutPattern:   true,
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc(loginPattern, s.login)
	mux.HandleFunc(callbackPattern, s.callback)
	mux.HandleFunc(logoutPattern, s.logout)
	mux.HandleFunc("GET /api/v1/me", me)
}

// Gate wraps the whole mux rather than single routes, so a path with no
// handler is refused like any other and a module added later starts protected.
//
// It asks the mux which pattern a request would reach and compares that with
// the public list. A test on the path instead would open any wildcard route
// that a public prefix happens to fit, such as /api/{module}/summary.
func (s *Service) Gate(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); publicPatterns[pattern] {
			mux.ServeHTTP(w, r)
			return
		}
		// Owner data and its refusals must not outlive the session in a proxy or browser cache.
		w.Header().Set("Cache-Control", "no-store")
		id, ok := s.session(r)
		if !ok {
			unauthorized(w)
			return
		}
		mux.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), identityKey{}, id)))
	})
}

type identityKey struct{}

// IdentityFrom returns who is signed in, for handlers behind the gate.
func IdentityFrom(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey{}).(Identity)
	return id, ok
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = io.WriteString(w, `{"error":"unauthorized"}`+"\n")
}

func me(w http.ResponseWriter, r *http.Request) {
	id, ok := IdentityFrom(r.Context())
	if !ok {
		unauthorized(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Provider string `json:"provider"`
		Name     string `json:"name"`
	}{id.Provider, id.Name})
}

func (s *Service) logout(w http.ResponseWriter, _ *http.Request) {
	s.clearCookie(w, sessionCookie)
	w.WriteHeader(http.StatusNoContent)
}
