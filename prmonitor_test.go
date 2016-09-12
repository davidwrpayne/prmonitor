package prmonitor

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Basic Auth Tests
func TestBasicAuthFailure(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic 490hf")

	BasicAuth("user", "pass", func(w http.ResponseWriter, r *http.Request) {

	})(w, r)

	if w.Code != 401 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 401, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("WWW-Authenticate") != "Basic" {
		t.Logf("ERROR: expected WWW-Authenticate header")
		t.Fail()
		return
	}
}

func TestBasicAuthSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	BasicAuth("user", "pass", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
		return
	})(w, r)

	if w.Code != 299 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 299, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("WWW-Authenticate") != "" {
		t.Logf("ERROR: unexpected WWW-Authenticate header on successful auth")
		t.Fail()
		return
	}
}

func TestBasicAuthSuccess2(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic Zm9vOmJhcg==")

	BasicAuth("foo", "bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
		return
	})(w, r)

	if w.Code != 299 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 299, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("WWW-Authenticate") != "" {
		t.Logf("ERROR: unexpected WWW-Authenticate header on successful auth")
		t.Fail()
		return
	}
}

// SSL Required Tests
func TestSSLRequiredRedirects(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://example.org/unsecure", nil)
	r.Header.Set("X-Forwarded-Proto", "http")

	SSLRequired("https://example.org/secure", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
		return
	})(w, r)

	if w.Code != 301 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 301, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("Location") != "https://example.org/secure" {
		t.Logf("ERROR: location '%s' expected, but got '%s'", "https://example.org/secure", w.Header().Get("Location"))
		t.Fail()
		return
	}
}

func TestSSLRequiredRedirects2(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://other.example.org/1", nil)
	r.Header.Set("X-Forwarded-Proto", "http")

	SSLRequired("https://other.example.org/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
		return
	})(w, r)

	if w.Code != 301 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 301, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("Location") != "https://other.example.org/1" {
		t.Logf("ERROR: location '%s' expected, but got '%s'", "https://other.example.org/1", w.Header().Get("Location"))
		t.Fail()
		return
	}
}

func TestSSLNotRedirectedIfAlreadyHTTPS(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "http://other.example.org/1", nil)
	r.Header.Set("X-Forwarded-Proto", "https")

	SSLRequired("http://other.example.org/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(299)
		return
	})(w, r)

	if w.Code != 299 {
		t.Logf("ERROR: http code '%d' expected, but got '%d'", 299, w.Code)
		t.Fail()
		return
	}

	if w.Header().Get("Location") != "" {
		t.Logf("ERROR: no location expected but got '%s'", w.Header().Get("Location"))
		t.Fail()
		return
	}
}

// Render Tests - attempt to render the page to a tmp file so it
// can be visually inspected.
func TestRender(t *testing.T) {
	prs := []SummarizedPullRequest{
		{
			Owner:    "brentdrich",
			Repo:     "prmonitor",
			Number:   4,
			Title:    "test pr",
			Author:   "brentdrich",
			OpenedAt: time.Now().Add(-5 * time.Hour),
		},
		{
			Owner:    "brentdrich",
			Repo:     "prmonitor",
			Number:   5,
			Title:    "yellow zone pr",
			Author:   "brentdrich",
			OpenedAt: time.Now().Add(-25 * time.Hour),
		},
		{
			Owner:    "brentdrich",
			Repo:     "prmonitor",
			Number:   6,
			Title:    "red zone pr",
			Author:   "brentdrich",
			OpenedAt: time.Now().Add(-73 * time.Hour),
		},
		{
			Owner:    "brentdrich",
			Repo:     "prmonitor",
			Number:   7,
			Title:    "boundary value pr",
			Author:   "brentdrich",
			OpenedAt: time.Now().Add(-1000 * time.Hour),
		},
	}
	f, err := os.Create("tmp.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	Render(f, prs)
}
