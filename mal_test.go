package mal

import (
	"fmt"
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

func testBasicAuth(t *testing.T, r *http.Request, usedWant bool, unameWant, passWant string) {
	uname, pass, used := r.BasicAuth()
	if used != usedWant {
		t.Errorf("BasicAuth used = %v, want %v", used, usedWant)
	}
	if uname != unameWant {
		t.Errorf("BasicAuth username = %v, want %v", uname, unameWant)
	}
	if pass != passWant {
		t.Errorf("BasicAuth password = %v, want %v", pass, passWant)
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

func TestClient_NewClient(t *testing.T) {
	c := NewClient()

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient.BaseURL = %v, want %v", got, want)
	}

	if got, want := c.UserAgent, defaultUserAgent; got != want {
		t.Errorf("NewClient.UserAgent = %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient()

	inURL, outURL := "/foo", defaultBaseURL+"foo"

	inData := &User{ID: 5, Username: "TestUser"}
	v := url.Values{}
	v.Set("data", "<user><id>5</id><username>TestUser</username></user>")
	urlEncData, _ := url.Parse(v.Encode())
	outData := urlEncData.Path

	req, _ := c.NewRequest("GET", inURL, inData)

	// test that the endpoint URL was correctly added to the base URL
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL =  %v, want %v", inURL, got, want)
	}

	// test that body was encoded to XML and then URL enconded as data=...
	body, _ := ioutil.ReadAll(req.Body)
	urlEncBody, _ := url.Parse(string(body)) // url.Path holds the URL decoded string
	if got, want := urlEncBody.Path, outData; got != want {
		t.Errorf("NewRequest(%+v) Body = '%v', want '%v'", inData, got, want)
	}

	testBasicAuth(t, req, false, "", "")
	testUserAgent(t, req, defaultUserAgent)

}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar string `xml:"bar"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if want := "GET"; r.Method != want {
			t.Errorf("request method = %v, want %v", r.Method, want)
		}
		fmt.Fprint(w, `<foo><bar>&bull; foobar</bar></foo>`)
	})

	req, _ := client.NewRequest("GET", "/", nil)

	body := new(foo)
	response, _ := client.Do(req, body)

	want := &foo{"&bull; foobar"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Do() response body = %v, want %v", body, want)
	}

	if want, got := "<foo><bar>&bull; foobar</bar></foo>", string(response.Body); want != got {
		t.Errorf("Do() mal.Response.Body = %v, want %v", got, want)
	}
}

func TestClient_Do_404(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)

	response, err := client.Do(req, nil)

	if err == nil {
		t.Errorf("expected 404 err")
	}

	// http.Error seems to be adding a new line to the message.
	if want, got := "not found\n", string(response.Body); want != got {
		t.Errorf("Do() mal.Response.Body = %v, want %v", got, want)
	}
}

// consider deleting
func TestClient_NewRequest_bad_endpoint(t *testing.T) {
	c := NewClient()
	inURL := "%foo"
	_, err := c.NewRequest("GET", inURL, nil)
	if err == nil {
		t.Errorf("NewRequest(%q) should return parse err", inURL)
	}
}

// consider deleting
func TestClient_NewRequest_xml_encode_err(t *testing.T) {
	c := NewClient()
	in := func() {}
	_, err := c.NewRequest("GET", "/foo", in)
	if err == nil {
		t.Errorf("NewRequest(%q) should return XML decode err", in)
	}
}
