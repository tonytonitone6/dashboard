package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"dashboard/api/internal/auth"
)

const (
	clientID      = "dashboard-client-id"
	clientSecret  = "dashboard-client-secret"
	sessionCookie = "dashboard_session"
	stateCookie   = "dashboard_oauth_state"
	loginPath     = "/api/auth/google/login"
	callbackPath  = "/api/auth/google/callback"
	mePath        = "/api/v1/me"
	thirtyDays    = 30 * 24 * time.Hour
)

var (
	signingKey = strings.Repeat("k", 32)
	otherKey   = strings.Repeat("o", 32)

	owner    = user{sub: "1001", name: "Stan", email: "stan@example.com"}
	stranger = user{sub: "2002", name: "Mallory", email: "mallory@example.com"}
)

type user struct{ sub, name, email string }

func (u user) code() string  { return "code-" + u.sub }
func (u user) token() string { return "token-" + u.sub }

// provider stands in for Google. Like the real one, it refuses an authorize or
// token request whose client, redirect URI or scope is not what was registered.
type provider struct {
	*httptest.Server
	t           *testing.T
	redirectURI string
	user        user
	cancels     bool
	// tokenMovedTo makes the token endpoint answer with a redirect to this URL.
	tokenMovedTo string
	// padding is the size of a filler field that precedes "sub" in the userinfo reply.
	padding int
}

func newProvider(t *testing.T) *provider {
	p := &provider{t: t, user: owner}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /authorize", p.authorize)
	mux.HandleFunc("POST /token", p.token)
	mux.HandleFunc("GET /userinfo", p.userinfo)
	p.Server = httptest.NewServer(mux)
	t.Cleanup(p.Close)
	return p
}

func (p *provider) expect(endpoint string, got url.Values, want map[string]string) bool {
	ok := true
	for key, value := range want {
		if got.Get(key) != value {
			p.t.Errorf("%s request has %s=%q, want %q", endpoint, key, got.Get(key), value)
			ok = false
		}
	}
	return ok
}

func (p *provider) authorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ok := p.expect("authorize", query, map[string]string{
		"client_id":     clientID,
		"redirect_uri":  p.redirectURI,
		"response_type": "code",
		"scope":         "openid profile",
		"prompt":        "select_account",
	})
	if query.Get("state") == "" {
		p.t.Error("authorize request has no state")
		ok = false
	}
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	back := url.Values{"state": {query.Get("state")}}
	if p.cancels {
		back.Set("error", "access_denied")
	} else {
		back.Set("code", p.user.code())
	}
	http.Redirect(w, r, p.redirectURI+"?"+back.Encode(), http.StatusFound)
}

func (p *provider) token(w http.ResponseWriter, r *http.Request) {
	if p.tokenMovedTo != "" {
		http.Redirect(w, r, p.tokenMovedTo, http.StatusTemporaryRedirect)
		return
	}
	_ = r.ParseForm()
	ok := p.expect("token", r.PostForm, map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  p.redirectURI,
	})
	if r.PostForm.Get("code") == "" {
		p.t.Error("token request has no code to exchange")
	}
	w.Header().Set("Content-Type", "application/json")
	if !ok || r.PostForm.Get("code") != p.user.code() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token": p.user.token(),
		"token_type":   "Bearer",
		"expires_in":   3599,
	})
}

func (p *provider) userinfo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+p.user.token() {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"sub":   p.user.sub,
		"name":  p.user.name,
		"email": p.user.email,
		"pad":   strings.Repeat("x", p.padding),
	})
}

// site is the API as main.go assembles it: the gate around a mux that holds
// health, the auth routes and one module.
type site struct {
	t       *testing.T
	server  *httptest.Server
	origin  string
	google  *provider
	browser *http.Client
	now     time.Time
}

type options struct {
	tls   bool
	key   string
	owner string
	// alias registers a second provider under this name that reaches the same fake.
	alias string
}

func newSite(t *testing.T, o options) *site {
	t.Helper()
	if o.key == "" {
		o.key = signingKey
	}
	if o.owner == "" {
		o.owner = owner.sub
	}
	s := &site{
		t:      t,
		server: httptest.NewUnstartedServer(nil),
		google: newProvider(t),
		now:    time.Date(2026, time.November, 1, 9, 0, 0, 0, time.UTC),
	}
	scheme := "http"
	if o.tls {
		scheme = "https"
	}
	s.origin = scheme + "://" + s.server.Listener.Addr().String()
	s.google.redirectURI = s.origin + callbackPath

	cfg, err := auth.ConfigFromEnv(lookup(map[string]string{
		"PUBLIC_URL":           s.origin,
		"SESSION_KEY":          o.key,
		"GOOGLE_CLIENT_ID":     clientID,
		"GOOGLE_CLIENT_SECRET": clientSecret,
		"OWNER_GOOGLE_SUB":     o.owner,
	}))
	if err != nil {
		t.Fatalf("test configuration: %v", err)
	}
	for name, entry := range cfg.Providers {
		entry.AuthorizeURL = s.google.URL + "/authorize"
		entry.TokenURL = s.google.URL + "/token"
		entry.UserinfoURL = s.google.URL + "/userinfo"
		cfg.Providers[name] = entry
	}
	if o.alias != "" {
		cfg.Providers[o.alias] = cfg.Providers["google"]
	}
	service := auth.New(cfg, func() time.Time { return s.now })

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "ok")
	})
	service.Register(mux)
	mux.HandleFunc("GET /api/v1/module", func(w http.ResponseWriter, r *http.Request) {
		id, _ := auth.IdentityFrom(r.Context())
		fmt.Fprintf(w, "%s %s %s", id.Provider, id.Subject, id.Name)
	})
	mux.HandleFunc("GET /api/{module}/summary", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "summary of %s", r.PathValue("module"))
	})
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "file %s", r.PathValue("path"))
	})
	s.server.Config.Handler = service.Gate(mux)
	if o.tls {
		s.server.StartTLS()
	} else {
		s.server.Start()
	}
	t.Cleanup(s.server.Close)
	s.browser = s.newBrowser()
	return s
}

func noRedirects(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

func (s *site) newBrowser() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		s.t.Fatal(err)
	}
	return &http.Client{Transport: s.server.Client().Transport, Jar: jar, CheckRedirect: noRedirects}
}

type reply struct {
	*http.Response
	body string
}

func (r reply) cookie(name string) *http.Cookie {
	for _, c := range r.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (s *site) send(browser *http.Client, method, target string, edit ...func(*http.Request)) reply {
	s.t.Helper()
	if strings.HasPrefix(target, "/") {
		target = s.origin + target
	}
	req, err := http.NewRequest(method, target, nil)
	if err != nil {
		s.t.Fatal(err)
	}
	for _, e := range edit {
		e(req)
	}
	resp, err := browser.Do(req)
	if err != nil {
		s.t.Fatalf("%s %s: %v", method, target, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.t.Fatalf("%s %s: reading the body: %v", method, target, err)
	}
	return reply{resp, string(body)}
}

func (s *site) get(target string) reply {
	s.t.Helper()
	return s.send(s.browser, http.MethodGet, target)
}

// getWith presents one session cookie from a browser that holds nothing else.
func (s *site) getWith(session, path string) reply {
	s.t.Helper()
	return s.send(s.newBrowser(), http.MethodGet, path, func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: sessionCookie, Value: session})
	})
}

// approve visits the consent screen a login reply points at and returns the
// URL the provider sends the browser back to.
func (s *site) approve(login reply) string {
	s.t.Helper()
	consent := s.send(&http.Client{CheckRedirect: noRedirects}, http.MethodGet, login.Header.Get("Location"))
	if consent.StatusCode != http.StatusFound {
		s.t.Fatalf("provider answered the authorize request with %d: %s", consent.StatusCode, consent.body)
	}
	return consent.Header.Get("Location")
}

func (s *site) signIn(as user) reply {
	s.t.Helper()
	s.google.user = as
	return s.get(s.approve(s.get(loginPath)))
}

func (s *site) sessionOf(as user) string {
	s.t.Helper()
	c := s.signIn(as).cookie(sessionCookie)
	if c == nil {
		s.t.Fatalf("signing in as %s set no session cookie", as.name)
	}
	return c.Value
}

func editQuery(target string, edit func(url.Values)) string {
	u, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	query := u.Query()
	edit(query)
	u.RawQuery = query.Encode()
	return u.String()
}

func captureLog(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &buf
}

func TestPublicRoutesNeedNoSession(t *testing.T) {
	s := newSite(t, options{})
	for _, tc := range []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/api/health", http.StatusOK},
		{http.MethodHead, "/api/health", http.StatusOK},
		{http.MethodGet, loginPath, http.StatusFound},
		{http.MethodPost, "/api/auth/logout", http.StatusNoContent},
		{http.MethodGet, "/api/auth/nobody/login", http.StatusNotFound},
		{http.MethodGet, "/api/auth/nobody/callback", http.StatusNotFound},
	} {
		if got := s.send(s.newBrowser(), tc.method, tc.path).StatusCode; got != tc.want {
			t.Errorf("%s %s without a session = %d, want %d", tc.method, tc.path, got, tc.want)
		}
	}
	if body := s.get("/api/health").body; body != "ok" {
		t.Errorf("health body = %q, want ok", body)
	}
}

func TestEverythingElseNeedsASession(t *testing.T) {
	s := newSite(t, options{})
	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, mePath},
		{http.MethodGet, "/api/v1/module"},
		{http.MethodGet, "/api/v1/does-not-exist"},
		{http.MethodGet, "/"},
		{http.MethodPost, "/api/health"},
		{http.MethodGet, "/api/health/deeper"},
		{http.MethodGet, "/api/healthz"},
		{http.MethodGet, "/api/authority/google/login"},
		{http.MethodGet, "/api/auth/no-such-route"},
		{http.MethodGet, "/api/v1/summary"},
		{http.MethodGet, "/api/auth/summary"},
		{http.MethodGet, "/api/auth%2Fx/summary"},
		{http.MethodGet, "/api/auth%2F..%2Fv1/summary"},
		{http.MethodGet, "/files/api/auth/google/login"},
		{http.MethodGet, "/api/auth/../v1/module"},
	} {
		got := s.send(s.newBrowser(), tc.method, tc.path)
		if got.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s %s without a session = %d, want 401", tc.method, tc.path, got.StatusCode)
			continue
		}
		if cache := got.Header.Get("Cache-Control"); cache != "no-store" {
			t.Errorf("%s %s: 401 has Cache-Control %q, want no-store", tc.method, tc.path, cache)
		}
		if contentType := got.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("%s %s: 401 has Content-Type %q, want application/json", tc.method, tc.path, contentType)
		}
		if !json.Valid([]byte(got.body)) {
			t.Errorf("%s %s: 401 body is not JSON: %q", tc.method, tc.path, got.body)
		}
	}
}

func TestOwnerSignsInAndReadsGatedRoutes(t *testing.T) {
	s := newSite(t, options{})

	callback := s.signIn(owner)
	if callback.StatusCode != http.StatusFound || callback.Header.Get("Location") != "/" {
		t.Fatalf("owner's callback = %d to %q, want 302 to /", callback.StatusCode, callback.Header.Get("Location"))
	}

	me := s.get(mePath)
	if me.StatusCode != http.StatusOK || me.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("%s after sign-in = %d %q, want 200 application/json", mePath, me.StatusCode, me.Header.Get("Content-Type"))
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(me.body), &got); err != nil {
		t.Fatalf("%s body %q: %v", mePath, me.body, err)
	}
	if want := map[string]string{"provider": "google", "name": owner.name}; !maps.Equal(got, want) {
		t.Errorf("%s = %v, want %v", mePath, got, want)
	}

	module := s.get("/api/v1/module")
	if want := "google 1001 Stan"; module.StatusCode != http.StatusOK || module.body != want {
		t.Errorf("module route saw %d %q, want 200 %q", module.StatusCode, module.body, want)
	}
}

func TestMeRefusesWhenTheMuxIsServedWithoutTheGate(t *testing.T) {
	cfg, err := auth.ConfigFromEnv(lookup(completeEnv()))
	if err != nil {
		t.Fatalf("test configuration: %v", err)
	}
	mux := http.NewServeMux()
	auth.New(cfg, time.Now).Register(mux)

	got := httptest.NewRecorder()
	mux.ServeHTTP(got, httptest.NewRequest(http.MethodGet, mePath, nil))
	if got.Code != http.StatusUnauthorized {
		t.Errorf("%s on an ungated mux = %d, want 401", mePath, got.Code)
	}
}

func TestLoginRedirectIgnoresTheHostHeader(t *testing.T) {
	s := newSite(t, options{})
	states := map[string]bool{}
	for range 2 {
		login := s.send(s.browser, http.MethodGet, loginPath, func(r *http.Request) { r.Host = "attacker.example" })
		location, err := url.Parse(login.Header.Get("Location"))
		if err != nil || login.StatusCode != http.StatusFound {
			t.Fatalf("login = %d to %q, want a 302", login.StatusCode, login.Header.Get("Location"))
		}
		if got, want := location.Scheme+"://"+location.Host+location.Path, s.google.URL+"/authorize"; got != want {
			t.Errorf("login redirects to %s, want %s", got, want)
		}
		if got, want := location.Query().Get("redirect_uri"), s.origin+callbackPath; got != want {
			t.Errorf("redirect_uri = %s, want %s", got, want)
		}
		states[location.Query().Get("state")] = true
	}
	if len(states) != 2 {
		t.Errorf("two logins used the states %v, want two different ones", states)
	}
}

type attributes struct {
	path     string
	maxAge   int
	httpOnly bool
	secure   bool
	sameSite http.SameSite
}

func TestCookieAttributes(t *testing.T) {
	for _, tc := range []struct {
		name string
		tls  bool
	}{{"http site", false}, {"https site", true}} {
		t.Run(tc.name, func(t *testing.T) {
			s := newSite(t, options{tls: tc.tls})
			login := s.get(loginPath)
			callback := s.get(s.approve(login))

			check := func(from reply, name string, want attributes) {
				t.Helper()
				c := from.cookie(name)
				if c == nil {
					t.Errorf("no %s cookie was set", name)
					return
				}
				if got := (attributes{c.Path, c.MaxAge, c.HttpOnly, c.Secure, c.SameSite}); got != want {
					t.Errorf("%s cookie = %+v, want %+v", name, got, want)
				}
			}
			check(login, stateCookie, attributes{"/api/auth", 10 * 60, true, tc.tls, http.SameSiteLaxMode})
			check(callback, sessionCookie, attributes{"/", 30 * 24 * 60 * 60, true, tc.tls, http.SameSiteStrictMode})
		})
	}
}

func TestNonOwnerIsDenied(t *testing.T) {
	for _, tc := range []struct {
		name string
		user user
	}{
		{"a stranger", stranger},
		{"the owner's name and email under another ID", user{sub: "3003", name: owner.name, email: owner.email}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logs := captureLog(t)
			s := newSite(t, options{})

			callback := s.signIn(tc.user)
			if callback.StatusCode != http.StatusFound || callback.Header.Get("Location") != "/?login=denied" {
				t.Errorf("callback = %d to %q, want 302 to /?login=denied", callback.StatusCode, callback.Header.Get("Location"))
			}
			if callback.cookie(sessionCookie) != nil {
				t.Error("a session cookie was set")
			}
			if got := s.get(mePath).StatusCode; got != http.StatusUnauthorized {
				t.Errorf("%s after a denied sign-in = %d, want 401", mePath, got)
			}

			logged := false
			for line := range strings.Lines(logs.String()) {
				logged = logged || (strings.Contains(line, "google") && strings.Contains(line, tc.user.sub))
			}
			if !logged {
				t.Errorf("no log line names the provider and the refused ID %s:\n%s", tc.user.sub, logs)
			}
		})
	}
}

func TestFailedSignInGetsNoSession(t *testing.T) {
	for _, tc := range []struct {
		name     string
		callback func(s *site) reply
	}{
		{"the state differs from the cookie", func(s *site) reply {
			return s.get(editQuery(s.approve(s.get(loginPath)), func(q url.Values) { q.Set("state", "forged") }))
		}},
		{"the callback carries no state", func(s *site) reply {
			return s.get(editQuery(s.approve(s.get(loginPath)), func(q url.Values) { q.Del("state") }))
		}},
		{"the browser has no state cookie", func(s *site) reply {
			return s.send(s.newBrowser(), http.MethodGet, s.approve(s.get(loginPath)))
		}},
		{"an empty state cookie meets a callback with no state", func(s *site) reply {
			return s.send(s.newBrowser(), http.MethodGet, callbackPath+"?code="+owner.code(), func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: stateCookie, Value: ""})
			})
		}},
		{"the user cancels at the provider", func(s *site) reply {
			s.google.cancels = true
			return s.get(s.approve(s.get(loginPath)))
		}},
		{"the provider rejects the code", func(s *site) reply {
			return s.get(editQuery(s.approve(s.get(loginPath)), func(q url.Values) { q.Set("code", "forged") }))
		}},
		{"the provider is unreachable during the code exchange", func(s *site) reply {
			back := s.approve(s.get(loginPath))
			s.google.Close()
			return s.get(back)
		}},
		{"the provider reports no user ID", func(s *site) reply {
			return s.signIn(user{name: "Nobody"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newSite(t, options{})
			callback := tc.callback(s)
			if callback.StatusCode != http.StatusFound || callback.Header.Get("Location") != "/?login=failed" {
				t.Errorf("callback = %d to %q, want 302 to /?login=failed", callback.StatusCode, callback.Header.Get("Location"))
			}
			if callback.cookie(sessionCookie) != nil {
				t.Error("a session cookie was set")
			}
			if got := s.get(mePath).StatusCode; got != http.StatusUnauthorized {
				t.Errorf("%s after a failed sign-in = %d, want 401", mePath, got)
			}
		})
	}
}

func TestCallbackCannotBeReplayed(t *testing.T) {
	s := newSite(t, options{})
	back := s.approve(s.get(loginPath))
	if got := s.get(back).Header.Get("Location"); got != "/" {
		t.Fatalf("first callback redirects to %q, want /", got)
	}
	if got := s.get(back).Header.Get("Location"); got != "/?login=failed" {
		t.Errorf("replayed callback redirects to %q, want /?login=failed", got)
	}
}

func TestSessionLastsThirtyDaysFromSignIn(t *testing.T) {
	s := newSite(t, options{})
	signedInAt := time.Date(2026, time.November, 1, 9, 0, 0, 0, time.UTC)
	s.now = signedInAt
	s.signIn(owner)

	for _, tc := range []struct {
		at   time.Time
		want int
	}{
		{time.Date(2026, time.November, 30, 9, 0, 0, 0, time.UTC), http.StatusOK},
		{signedInAt.Add(thirtyDays - time.Second), http.StatusOK},
		{signedInAt.Add(thirtyDays), http.StatusUnauthorized},
		{signedInAt.Add(thirtyDays + time.Minute), http.StatusUnauthorized},
		{time.Date(2026, time.December, 2, 9, 0, 0, 0, time.UTC), http.StatusUnauthorized},
	} {
		s.now = tc.at
		got := s.get(mePath)
		if got.StatusCode != tc.want {
			t.Errorf("request at %s = %d, want %d", tc.at.Format(time.DateTime), got.StatusCode, tc.want)
		}
		if len(got.Cookies()) != 0 {
			t.Errorf("request at %s set a cookie, so using the dashboard would renew the session", tc.at.Format(time.DateTime))
		}
	}
}

func TestAlteredCookieIsRefused(t *testing.T) {
	s := newSite(t, options{})
	session := s.sessionOf(owner)
	if got := s.getWith(session, mePath).StatusCode; got != http.StatusOK {
		t.Fatalf("untouched cookie = %d, want 200", got)
	}
	for i := range len(session) {
		changed := []byte(session)
		changed[i] = 'A'
		if session[i] == 'A' {
			changed[i] = 'B'
		}
		if got := s.getWith(string(changed), mePath).StatusCode; got != http.StatusUnauthorized {
			t.Fatalf("cookie with character %d changed = %d, want 401", i, got)
		}
		if got := s.getWith(session[:i], mePath).StatusCode; got != http.StatusUnauthorized {
			t.Fatalf("cookie cut to %d characters = %d, want 401", i, got)
		}
	}
}

func TestCookieSignedWithAnotherKeyIsRefused(t *testing.T) {
	previous := newSite(t, options{key: otherKey})
	current := newSite(t, options{})

	session := previous.sessionOf(owner)
	if got := previous.getWith(session, mePath).StatusCode; got != http.StatusOK {
		t.Fatalf("cookie on the site that signed it = %d, want 200", got)
	}
	if got := current.getWith(session, mePath).StatusCode; got != http.StatusUnauthorized {
		t.Errorf("cookie signed with another key = %d, want 401", got)
	}
}

func TestSessionEndsWhenTheOwnerChanges(t *testing.T) {
	session := newSite(t, options{}).sessionOf(owner)

	restarted := newSite(t, options{})
	if got := restarted.getWith(session, mePath).StatusCode; got != http.StatusOK {
		t.Fatalf("cookie after a restart with the same configuration = %d, want 200", got)
	}
	reconfigured := newSite(t, options{owner: stranger.sub})
	if got := reconfigured.getWith(session, mePath).StatusCode; got != http.StatusUnauthorized {
		t.Errorf("cookie after the owner ID changed = %d, want 401", got)
	}
}

func TestLogoutEndsTheSession(t *testing.T) {
	s := newSite(t, options{})
	s.signIn(owner)
	if got := s.get(mePath).StatusCode; got != http.StatusOK {
		t.Fatalf("%s after sign-in = %d, want 200", mePath, got)
	}
	if got := s.send(s.browser, http.MethodPost, "/api/auth/logout").StatusCode; got != http.StatusNoContent {
		t.Errorf("logout = %d, want 204", got)
	}
	if got := s.get(mePath).StatusCode; got != http.StatusUnauthorized {
		t.Errorf("%s after logout = %d, want 401", mePath, got)
	}
}

func TestLogCarriesNoSecrets(t *testing.T) {
	logs := captureLog(t)
	s := newSite(t, options{})

	session := s.sessionOf(owner)
	s.signIn(stranger)
	s.google.user = owner
	login := s.get(loginPath)
	state := login.cookie(stateCookie).Value
	s.get(editQuery(s.approve(login), func(q url.Values) { q.Set("code", "forged-code") }))
	s.get(editQuery(s.approve(s.get(loginPath)), func(q url.Values) { q.Set("state", "forged-state") }))
	back := s.approve(s.get(loginPath))
	s.google.Close()
	s.get(back)

	if got := strings.Count(logs.String(), "\n"); got < 4 {
		t.Fatalf("%d lines were logged, want at least one for each of the four refused sign-ins:\n%s", got, logs)
	}
	for what, secret := range map[string]string{
		"client secret":    clientSecret,
		"session key":      signingKey,
		"owner's code":     owner.code(),
		"stranger's code":  stranger.code(),
		"forged code":      "forged-code",
		"owner's token":    owner.token(),
		"stranger's token": stranger.token(),
		"session cookie":   session,
		"state":            state,
		"forged state":     "forged-state",
	} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("the log contains the %s:\n%s", what, logs)
		}
	}
}

func TestOwnerDataIsNotCacheable(t *testing.T) {
	s := newSite(t, options{})
	s.signIn(owner)
	for _, path := range []string{mePath, "/api/v1/module"} {
		got := s.get(path)
		if got.StatusCode != http.StatusOK {
			t.Fatalf("GET %s as the owner = %d, want 200", path, got.StatusCode)
		}
		if cache := got.Header.Get("Cache-Control"); cache != "no-store" {
			t.Errorf("GET %s as the owner has Cache-Control %q, want no-store", path, cache)
		}
	}
}

func TestOwnerOfOneProviderIsNotOwnerOfAnother(t *testing.T) {
	s := newSite(t, options{alias: "other"})
	s.google.redirectURI = s.origin + "/api/auth/other/callback"
	s.google.user = owner

	done := s.get(s.approve(s.get("/api/auth/other/login")))
	if got := done.Header.Get("Location"); got != "/?login=denied" {
		t.Errorf("the owner's ID arriving from another provider redirects to %q, want /?login=denied", got)
	}
	if done.cookie(sessionCookie) != nil {
		t.Error("the owner's ID arriving from another provider got a session")
	}
}

func TestProviderRedirectIsNotFollowed(t *testing.T) {
	s := newSite(t, options{})
	var leaked []string
	elsewhere := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		leaked = append(leaked, r.Method+" "+r.URL.Path+" "+r.PostForm.Encode())
	}))
	t.Cleanup(elsewhere.Close)
	s.google.tokenMovedTo = elsewhere.URL + "/token"

	done := s.signIn(owner)
	if got := done.Header.Get("Location"); got != "/?login=failed" {
		t.Errorf("sign-in with a redirecting token endpoint goes to %q, want /?login=failed", got)
	}
	if done.cookie(sessionCookie) != nil {
		t.Error("sign-in with a redirecting token endpoint got a session")
	}
	if len(leaked) != 0 {
		t.Errorf("the code exchange followed the redirect and sent: %v", leaked)
	}
}

func TestOversizedProviderReplyFailsSignIn(t *testing.T) {
	s := newSite(t, options{})
	s.google.padding = 1 << 17

	done := s.signIn(owner)
	if got := done.Header.Get("Location"); got != "/?login=failed" {
		t.Errorf("sign-in with a 128 KiB userinfo reply goes to %q, want /?login=failed", got)
	}
	if done.cookie(sessionCookie) != nil {
		t.Error("sign-in with a 128 KiB userinfo reply got a session")
	}
}
