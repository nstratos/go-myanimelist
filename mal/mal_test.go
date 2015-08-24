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
	// client is the MyAnimeList client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server

	// mux is the HTTP request multiplexer that the test HTTP server will use
	// to mock API responses.
	mux *http.ServeMux
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

func TestNewClient(t *testing.T) {
	c := NewClient()

	// test default base URL
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient.BaseURL = %v, want %v", got, want)
	}

	// test default user agent
	if got, want := c.UserAgent, defaultUserAgent; got != want {
		t.Errorf("NewClient.UserAgent = %v, want %v", got, want)
	}

	// test account default endpoint
	if got, want := c.Account.Endpoint.String(), defaultAccountEndpoint; got != want {
		t.Errorf("NewClient.Account.Endpoint = %v, want %v", got, want)
	}

	// test anime default endpoints
	if got, want := c.Anime.AddEndpoint.String(), defaultAnimeAddEndpoint; got != want {
		t.Errorf("NewClient.Anime.Addpoint = %v, want %v", got, want)
	}
	if got, want := c.Anime.UpdateEndpoint.String(), defaultAnimeUpdateEndpoint; got != want {
		t.Errorf("NewClient.Anime.UpdateEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Anime.DeleteEndpoint.String(), defaultAnimeDeleteEndpoint; got != want {
		t.Errorf("NewClient.Anime.DeleteEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Anime.SearchEndpoint.String(), defaultAnimeSearchEndpoint; got != want {
		t.Errorf("NewClient.Anime.SearchEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Anime.ListEndpoint.String(), defaultListEndpoint; got != want {
		t.Errorf("NewClient.Anime.ListEndpoint = %v, want %v", got, want)
	}

	// test manga default endpoints
	if got, want := c.Manga.AddEndpoint.String(), defaultMangaAddEndpoint; got != want {
		t.Errorf("NewClient.Manga.Addpoint = %v, want %v", got, want)
	}
	if got, want := c.Manga.UpdateEndpoint.String(), defaultMangaUpdateEndpoint; got != want {
		t.Errorf("NewClient.Manga.UpdateEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Manga.DeleteEndpoint.String(), defaultMangaDeleteEndpoint; got != want {
		t.Errorf("NewClient.Manga.DeleteEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Manga.SearchEndpoint.String(), defaultMangaSearchEndpoint; got != want {
		t.Errorf("NewClient.Manga.SearchEndpoint = %v, want %v", got, want)
	}
	if got, want := c.Manga.ListEndpoint.String(), defaultListEndpoint; got != want {
		t.Errorf("NewClient.Manga.ListEndpoint = %v, want %v", got, want)
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

func TestClient_Do_invalidXMLEntity(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar string `xml:"bar"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if want := "GET"; r.Method != want {
			t.Errorf("request method = %v, want %v", r.Method, want)
		}
		fmt.Fprint(w, `<foo><bar>&foo; bar</bar></foo>`)
	})

	req, _ := client.NewRequest("GET", "/", nil)

	body := new(foo)
	response, err := client.Do(req, body)

	if err == nil {
		t.Errorf("Do() receiving XML with invalid entity should return err")
	}

	if response == nil {
		t.Errorf("Do() receiving XML with invalid entity should also return response")
	}

	if want, got := "<foo><bar>&foo; bar</bar></foo>", string(response.Body); want != got {
		t.Errorf("Do() receiving XML with invalid entity mal.Response.Body = %v, want %v", got, want)
	}
}

func TestClient_Do_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)

	response, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	// http.Error seems to be adding a new line to the message.
	if want, got := "not found\n", string(response.Body); want != got {
		t.Errorf("Do() mal.Response.Body = %v, want %v", got, want)
	}
}

func TestClient_Do_connectionRefused(t *testing.T) {
	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)
	if err == nil {
		t.Error("Expected connection refused error.")
	}
}

func TestClient_post_invalidID(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/update/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testID(t, r, "0")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		// zeroEntry defined in anime_test.go
		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, "onhold"))
		http.Error(w, "Invalid ID", http.StatusNotImplemented)
	})

	response, err := client.post("api/animelist/update/", 0, AnimeEntry{Status: "onhold"})

	if err == nil {
		t.Errorf("Anime.Update invalid ID should return err")
	}

	if response == nil {
		t.Errorf("Anime.Update invalid ID should return also return response")
	}
}

func TestClient_delete_invalidID(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/delete/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testID(t, r, "0")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		http.Error(w, "Invalid ID", http.StatusNotImplemented)
	})

	response, err := client.delete("api/animelist/delete/", 0)

	if err == nil {
		t.Errorf("Anime.Delete invalid ID should return err")
	}

	if response == nil {
		t.Errorf("Anime.Delete invalid ID should return also return response")
	}
}

func TestClient_NewRequest_badEndpoint(t *testing.T) {
	c := NewClient()
	inURL := "%foo"
	_, err := c.NewRequest("GET", inURL, nil)
	if err == nil {
		t.Errorf("NewRequest(%q) should return parse err", inURL)
	}
}

func TestClient_NewRequest_xmlEncodeError(t *testing.T) {
	c := NewClient()
	in := func() {} // xml.Marshal cannot encode a func
	_, err := c.NewRequest("GET", "/foo", in)
	if err == nil {
		t.Errorf("NewRequest receiving a function as body should return XML encode err")
	}
}
