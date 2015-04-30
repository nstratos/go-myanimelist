package mal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the MyAnimeList client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a mal.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// mal client configured to use test server
	client = NewClient()
	client.BaseURL, _ = url.Parse(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testBasicAuth(t *testing.T, r *http.Request, uwant, pwant string) {
	uname, pass, ok := r.BasicAuth()
	if !ok {
		t.Errorf("BasicAuth used = %v, want %v", ok, true)
	}
	if uname != uwant {
		t.Errorf("BasicAuth username = %v, want %v", uname, uwant)
	}
	if pass != pwant {
		t.Errorf("BasicAuth password = %v, want %v", pass, pwant)
	}
}

func testUserAgent(t *testing.T, r *http.Request, want string) {
	agent := r.Header.Get("User-Agent")
	if want != agent {
		t.Errorf("User-Agent = %v, want %v", agent, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Unable to read body")
	}
	body := string(b)
	if body != want {
		t.Errorf("body = %v, want %v", body, want)
	}
}

func testContentType(t *testing.T, r *http.Request, want string) {
	ct := r.Header.Get("Content-Type")
	if ct != want {
		t.Errorf("Content-Type = %v, want %v", ct, want)
	}
}
