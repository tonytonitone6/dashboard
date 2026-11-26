package auth

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// Provider is one row of the sign-in table: everything that differs between
// OAuth providers. Adding a provider means adding a row in ConfigFromEnv.
type Provider struct {
	AuthorizeURL string
	TokenURL     string
	UserinfoURL  string
	Scope        string
	// AuthorizeExtra holds authorize parameters only this provider understands.
	AuthorizeExtra map[string]string
	ClientID       string
	ClientSecret   string
	// User reads the permanent user ID and the display name from a userinfo response.
	User func(userinfo []byte) (subject, name string, err error)
}

type owner struct{ provider, subject string }

// Config is the validated environment. The fields that guard access are
// unexported so that ConfigFromEnv is the only way to fill them.
type Config struct {
	// Providers is exported so tests can point the endpoint URLs at a local fake.
	Providers map[string]Provider

	publicURL  string
	secure     bool
	sessionKey []byte
	owners     map[owner]bool
}

func (c Config) isOwner(id Identity) bool {
	return c.owners[owner{id.Provider, id.Subject}]
}

const minSessionKeyBytes = 32

// ConfigFromEnv reports every problem in one error, naming variables and never
// their values.
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	var missing, padded, problems []string
	required := func(name string) string {
		value := getenv(name)
		switch {
		case value == "":
			missing = append(missing, name)
		case value != strings.TrimSpace(value):
			padded = append(padded, name)
		}
		return value
	}
	publicURL := required("PUBLIC_URL")
	sessionKey := required("SESSION_KEY")
	googleClientID := required("GOOGLE_CLIENT_ID")
	googleClientSecret := required("GOOGLE_CLIENT_SECRET")
	googleOwner := required("OWNER_GOOGLE_SUB")

	if len(missing) > 0 {
		problems = append(problems, "missing "+strings.Join(missing, ", "))
	}
	if len(padded) > 0 {
		problems = append(problems, "leading or trailing whitespace in "+strings.Join(padded, ", "))
	}
	if publicURL != "" && !isOrigin(publicURL) {
		problems = append(problems, "PUBLIC_URL must be an http or https origin such as https://example.com, with no trailing slash or path")
	}
	if sessionKey != "" && len(sessionKey) < minSessionKeyBytes {
		problems = append(problems, "SESSION_KEY must be at least 32 bytes")
	}
	if len(problems) > 0 {
		return Config{}, errors.New("configuration: " + strings.Join(problems, "; "))
	}

	return Config{
		Providers: map[string]Provider{
			"google": {
				AuthorizeURL: "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:     "https://oauth2.googleapis.com/token",
				UserinfoURL:  "https://openidconnect.googleapis.com/v1/userinfo",
				Scope:        "openid profile",
				// Without this Google silently reuses the browser's one signed-in account, so a
				// refused account could never be swapped for the owner's.
				AuthorizeExtra: map[string]string{"prompt": "select_account"},
				ClientID:       googleClientID,
				ClientSecret:   googleClientSecret,
				User:           googleUser,
			},
		},
		publicURL:  publicURL,
		secure:     strings.HasPrefix(publicURL, "https://"),
		sessionKey: []byte(sessionKey),
		owners:     map[owner]bool{{"google", googleOwner}: true},
	}, nil
}

// A path can never work: the callback URL, the cookie paths and the redirects
// after sign-in all assume the site is served from the root.
//
// The host checks reject spellings that parse but that a provider would not
// accept as the registered callback URL.
func isOrigin(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || raw != u.Scheme+"://"+u.Host {
		return false
	}
	host := u.Hostname()
	if host == "" || host != strings.ToLower(host) || strings.HasPrefix(host, ".") || strings.HasSuffix(host, ".") {
		return false
	}
	if strings.HasSuffix(u.Host, ":") {
		return false
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		return err == nil && n >= 1 && n <= 65535
	}
	return true
}

func googleUser(userinfo []byte) (subject, name string, err error) {
	var user struct {
		Sub  string `json:"sub"`
		Name string `json:"name"`
	}
	err = json.Unmarshal(userinfo, &user)
	return user.Sub, user.Name, err
}
