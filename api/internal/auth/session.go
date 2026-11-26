package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type cookie struct {
	name     string
	path     string
	sameSite http.SameSite
	lifetime time.Duration
}

var (
	// Strict is enough because the page shell is public and only the app's own fetch calls need the session.
	sessionCookie = cookie{"dashboard_session", "/", http.SameSiteStrictMode, 30 * 24 * time.Hour}
	// Lax, so the cookie survives the redirect back from the provider.
	stateCookie = cookie{"dashboard_oauth_state", "/api/auth", http.SameSiteLaxMode, 10 * time.Minute}
)

func (s *Service) setCookie(w http.ResponseWriter, c cookie, value string) {
	s.writeCookie(w, c, value, int(c.lifetime/time.Second))
}

func (s *Service) clearCookie(w http.ResponseWriter, c cookie) {
	s.writeCookie(w, c, "", -1)
}

func (s *Service) writeCookie(w http.ResponseWriter, c cookie, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.name,
		Value:    value,
		Path:     c.path,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   s.cfg.secure,
		SameSite: c.sameSite,
	})
}

type sessionPayload struct {
	Provider string `json:"provider"`
	Subject  string `json:"subject"`
	Name     string `json:"name"`
	Expires  int64  `json:"expires"`
}

// seal fixes the expiry at sign-in. Nothing re-issues the cookie later, so
// using the dashboard never extends a session.
func (s *Service) seal(id Identity) string {
	// Marshalling a struct of strings and an integer cannot fail.
	raw, _ := json.Marshal(sessionPayload{id.Provider, id.Subject, id.Name, s.now().Add(sessionCookie.lifetime).Unix()})
	payload := base64.RawURLEncoding.EncodeToString(raw)
	return payload + "." + s.sign(payload)
}

// sign covers the encoded payload exactly as it travels, and open compares the
// encoded signature, so no second spelling of a valid cookie exists.
func (s *Service) sign(payload string) string {
	mac := hmac.New(sha256.New, s.cfg.sessionKey)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *Service) open(value string) (Identity, bool) {
	payload, signature, _ := strings.Cut(value, ".")
	if !hmac.Equal([]byte(signature), []byte(s.sign(payload))) {
		return Identity{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Identity{}, false
	}
	var p sessionPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return Identity{}, false
	}
	if !s.now().Before(time.Unix(p.Expires, 0)) {
		return Identity{}, false
	}
	return Identity{Provider: p.Provider, Subject: p.Subject, Name: p.Name}, true
}

// session consults the owner list on every request, so removing an owner ends
// their sessions at the next restart.
func (s *Service) session(r *http.Request) (Identity, bool) {
	c, err := r.Cookie(sessionCookie.name)
	if err != nil {
		return Identity{}, false
	}
	id, ok := s.open(c.Value)
	return id, ok && s.cfg.isOwner(id)
}
