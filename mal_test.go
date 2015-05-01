package mal

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

func testID(t *testing.T, r *http.Request, want string) {
	idXML := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	xml := idXML[len(idXML)-4:]
	if xml != ".xml" {
		t.Errorf("URL path %v does not end in .xml", r.URL.Path)
	}
	id := idXML[:len(idXML)-4]
	if id != want {
		t.Errorf("provided id = %v, want %v", id, want)
	}
}

func testFormValue(t *testing.T, r *http.Request, value, want string) {
	v := r.FormValue(value)
	if v != want {
		t.Errorf("form value %v = %v, want %v", value, v, want)
	}
}

type urlValues map[string]string

func testURLValues(t *testing.T, r *http.Request, values urlValues) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}
	actual := r.URL.Query()
	if !reflect.DeepEqual(want, actual) {
		t.Errorf("URL Values = %v, want %v", actual, want)
	}
}
