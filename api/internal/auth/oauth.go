package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// outcome is how a callback ended. Its value is where the browser goes next.
type outcome string

const (
	signedIn outcome = "/"
	denied   outcome = "/?login=denied"
	failed   outcome = "/?login=failed"
)

const maxProviderResponseBytes = 1 << 16

// redirectURI is built from PUBLIC_URL alone. The request's Host header is
// whatever the client chose to send.
func (s *Service) redirectURI(provider string) string {
	return s.cfg.publicURL + "/api/auth/" + provider + "/callback"
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("provider")
	p, ok := s.cfg.Providers[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	state := rand.Text()
	s.setCookie(w, stateCookie, state)
	query := url.Values{
		"client_id":     {p.ClientID},
		"redirect_uri":  {s.redirectURI(name)},
		"response_type": {"code"},
		"scope":         {p.Scope},
		"state":         {state},
	}
	for key, value := range p.AuthorizeExtra {
		query.Set(key, value)
	}
	http.Redirect(w, r, p.AuthorizeURL+"?"+query.Encode(), http.StatusFound)
}

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("provider")
	p, ok := s.cfg.Providers[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	s.clearCookie(w, stateCookie)
	http.Redirect(w, r, string(s.finish(w, r, name, p)), http.StatusFound)
}

// finish sets the session cookie only on the signedIn outcome.
func (s *Service) finish(w http.ResponseWriter, r *http.Request, name string, p Provider) outcome {
	query := r.URL.Query()
	sent, err := r.Cookie(stateCookie.name)
	// The empty check stops a request with no state from matching an empty cookie.
	if err != nil || sent.Value == "" || subtle.ConstantTimeCompare([]byte(sent.Value), []byte(query.Get("state"))) != 1 {
		log.Printf("auth: %s sign-in failed: the state does not match this browser", name)
		return failed
	}
	code := query.Get("code")
	if code == "" {
		log.Printf("auth: %s sign-in failed: the provider sent no code", name)
		return failed
	}
	id, err := s.identify(r.Context(), name, p, code)
	if err != nil {
		log.Printf("auth: %s sign-in failed: %v", name, err)
		return failed
	}
	if !s.cfg.isOwner(id) {
		log.Printf("auth: sign-in refused, not an owner: provider=%s subject=%s", id.Provider, id.Subject)
		return denied
	}
	s.setCookie(w, sessionCookie, s.seal(id))
	return signedIn
}

// identify trades the code for a token and asks the provider who the token
// belongs to. Its errors never quote a response body, which can hold tokens.
func (s *Service) identify(ctx context.Context, name string, p Provider, code string) (Identity, error) {
	status, body, err := s.fetch(ctx, http.MethodPost, p.TokenURL, http.Header{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"application/json"},
	}, url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {s.redirectURI(name)},
		"client_id":     {p.ClientID},
		"client_secret": {p.ClientSecret},
	}.Encode())
	if err != nil {
		return Identity{}, fmt.Errorf("code exchange: %w", err)
	}
	var token struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	// A body that is not JSON leaves the token empty, which the next check reports.
	_ = json.Unmarshal(body, &token)
	if status != http.StatusOK || token.AccessToken == "" {
		return Identity{}, fmt.Errorf("code exchange: status %d, error %q", status, token.Error)
	}

	status, body, err = s.fetch(ctx, http.MethodGet, p.UserinfoURL, http.Header{
		"Authorization": {"Bearer " + token.AccessToken},
		"Accept":        {"application/json"},
	}, "")
	if err != nil {
		return Identity{}, fmt.Errorf("userinfo: %w", err)
	}
	if status != http.StatusOK {
		return Identity{}, fmt.Errorf("userinfo: status %d", status)
	}
	subject, display, err := p.User(body)
	if err != nil || subject == "" {
		return Identity{}, errors.New("userinfo: no user ID in the response")
	}
	return Identity{Provider: name, Subject: subject, Name: display}, nil
}

func (s *Service) fetch(ctx context.Context, method, target string, header http.Header, body string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, target, strings.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header = header
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	read, err := io.ReadAll(io.LimitReader(resp.Body, maxProviderResponseBytes))
	return resp.StatusCode, read, err
}
