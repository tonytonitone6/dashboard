package auth_test

import (
	"strings"
	"testing"

	"dashboard/api/internal/auth"
)

func completeEnv() map[string]string {
	return map[string]string{
		"PUBLIC_URL":           "https://dashboard.example",
		"SESSION_KEY":          signingKey,
		"GOOGLE_CLIENT_ID":     clientID,
		"GOOGLE_CLIENT_SECRET": clientSecret,
		"OWNER_GOOGLE_SUB":     owner.sub,
	}
}

func lookup(env map[string]string) func(string) string {
	return func(name string) string { return env[name] }
}

func TestConfigAcceptsACompleteEnvironment(t *testing.T) {
	for _, publicURL := range []string{"https://stanmao.com", "http://localhost:5173"} {
		env := completeEnv()
		env["PUBLIC_URL"] = publicURL
		if _, err := auth.ConfigFromEnv(lookup(env)); err != nil {
			t.Errorf("complete environment with PUBLIC_URL=%s: %v", publicURL, err)
		}
	}
}

func TestConfigNamesEachMissingVariable(t *testing.T) {
	for name := range completeEnv() {
		env := completeEnv()
		delete(env, name)
		_, err := auth.ConfigFromEnv(lookup(env))
		if err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("without %s: error %v, want one naming %s", name, err, name)
			continue
		}
		for other, value := range env {
			if strings.Contains(err.Error(), value) {
				t.Errorf("without %s: the error prints the value of %s", name, other)
			}
		}
	}
}

func TestConfigReportsEveryMissingVariableAtOnce(t *testing.T) {
	_, err := auth.ConfigFromEnv(lookup(nil))
	if err == nil {
		t.Fatal("an empty environment was accepted")
	}
	for name := range completeEnv() {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error for an empty environment does not name %s: %v", name, err)
		}
	}
}

func TestConfigRejectsAShortSessionKey(t *testing.T) {
	short := strings.Repeat("k", 31)
	env := completeEnv()
	env["SESSION_KEY"] = short
	_, err := auth.ConfigFromEnv(lookup(env))
	if err == nil || !strings.Contains(err.Error(), "SESSION_KEY") {
		t.Fatalf("31-byte SESSION_KEY: error %v, want one naming SESSION_KEY", err)
	}
	if strings.Contains(err.Error(), short) {
		t.Error("the error prints the session key")
	}
}

func TestConfigRejectsAPublicURLThatIsNotAnOrigin(t *testing.T) {
	for _, publicURL := range []string{
		"stanmao.com",
		"//stanmao.com",
		"ftp://stanmao.com",
		"http://",
		"https://stanmao.com/",
		"https://stanmao.com/dashboard",
		"https://stanmao.com?next=1",
		"https://stanmao.com:",
		"https://stanmao.com:99999",
		"https://stanmao.com:0",
		"https://STANMAO.com",
		"https://.",
		"https://stanmao.com.",
	} {
		env := completeEnv()
		env["PUBLIC_URL"] = publicURL
		_, err := auth.ConfigFromEnv(lookup(env))
		if err == nil || !strings.Contains(err.Error(), "PUBLIC_URL") {
			t.Errorf("PUBLIC_URL=%s: error %v, want one naming PUBLIC_URL", publicURL, err)
			continue
		}
		if strings.Contains(err.Error(), publicURL) {
			t.Errorf("PUBLIC_URL=%s: the error prints the value", publicURL)
		}
	}
}

func TestConfigRejectsSurroundingWhitespace(t *testing.T) {
	for name, value := range completeEnv() {
		for _, padded := range []string{" " + value, value + " ", value + "\t"} {
			env := completeEnv()
			env[name] = padded
			_, err := auth.ConfigFromEnv(lookup(env))
			if err == nil || !strings.Contains(err.Error(), name) {
				t.Errorf("%s with surrounding whitespace: error %v, want one naming %s", name, err, name)
			}
		}
	}
	env := completeEnv()
	env["SESSION_KEY"] = strings.Repeat(" ", 32)
	if _, err := auth.ConfigFromEnv(lookup(env)); err == nil {
		t.Error("a SESSION_KEY of 32 spaces was accepted")
	}
}
