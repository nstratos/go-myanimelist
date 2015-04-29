package mal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a mal.Client that is
// configured to talk to that test server.  Tests should register handlers on
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

func TestAccountService_Verify_credentials(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass", "TestAgent")

	mux.HandleFunc("/api/account/verify_credentials.xml", func(w http.ResponseWriter, r *http.Request) {
		testBasicAuth(t, r, "TestUser", "TestPass")
		// uname, pass, ok := r.BasicAuth()
		// if !ok || uname != "TestUser" || pass != "TestPass" {
		// 	http.Error(w, "credentials are wrong", http.StatusUnauthorized)
		// 	return
		// }
		testUserAgent(t, r, "TestAgent")
		// agent := r.Header.Get("User-Agent")
		// if agent != "TestAgent" {
		// 	http.Error(w, "user agent not recognized", http.StatusUnauthorized)
		// 	return
		// }
		testMethod(t, r, "GET")
		fmt.Fprint(w, `<user><id>1</id><username>TestUser</username></user>`)
	})

	user, _, err := client.Account.Verify() //should also return response!
	if err != nil {
		t.Errorf("Account.Verify returned error: %v", err)
	}

	want := &User{ID: 1, Username: "TestUser"}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Account.Verify returned %+v, want %+v", user, want)
	}
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
		t.Errorf("User Agent = %v, want %v", agent, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}
