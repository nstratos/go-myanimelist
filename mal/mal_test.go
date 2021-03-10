package mal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// setup sets up a test HTTP server along with a mal.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer that the test HTTP server will use
	// to mock API responses.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the MyAnimeList client being tested and is configured to use
	// test server.
	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL + "/")

	return client, mux, server.Close
}

type urlValues map[string]string

func testURLValues(t *testing.T, r *http.Request, values urlValues) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}
	actual := r.URL.Query()
	if !reflect.DeepEqual(want, actual) {
		t.Errorf("URL Values = %v, want %v", actual, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request body\nhave: %q\nwant: %q", got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func testContentType(t *testing.T, r *http.Request, want string) {
	ct := r.Header.Get("Content-Type")
	if ct != want {
		t.Errorf("Content-Type = %q, want %q", ct, want)
	}
}

func testErrorResponse(t *testing.T, err error, want ErrorResponse) {
	t.Helper()
	errResp := &ErrorResponse{}
	if !errors.As(err, &errResp) {
		t.Fatalf("err is type %T, want type *ErrorResponse.", err)
	}
	if got, want := errResp.Message, want.Message; got != want {
		t.Errorf("ErrorResponse.Message = %q, want %q", got, want)
	}
	if got, want := errResp.Err, want.Err; got != want {
		t.Errorf("ErrorResponse.Err = %q, want %q", got, want)
	}
}

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q)\nhave: %s\nwant: %s", v, j, w)
	}
}

func testResponseOffset(t *testing.T, resp *Response, next, prev int, prefix string) {
	t.Helper()
	if resp == nil {
		t.Fatalf("%s resp is nil, want NextOffset=%d and PrevOffset=%d", prefix, next, prev)
	}
	if got, want := resp.NextOffset, next; got != want {
		t.Errorf("%s resp.NextOffset=%d, want %d", prefix, got, want)
	}
	if got, want := resp.PrevOffset, prev; got != want {
		t.Errorf("%s resp.PrevOffset=%d, want %d", prefix, got, want)
	}
}

func testResponseStatusCode(t *testing.T, resp *Response, code int, prefix string) {
	t.Helper()
	if resp == nil {
		t.Fatalf("%s resp is nil, want StatusCode=%d", prefix, code)
	}
	if got, want := resp.StatusCode, code; got != want {
		t.Errorf("%s resp.StatusCode=%d, want %d", prefix, got, want)
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
	c := NewClient(nil)

	// test default base URL
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient.BaseURL = %v, want %v", got, want)
	}
}

// func TestNewClient_options(t *testing.T) {
// 	httpClient := &http.Client{}
// 	c := NewClient(
// 		Auth("TestUser", "TestPass"),
// 		HTTPClient(httpClient),
// 	)

// 	// test passing username and password as option
// 	if got, want := c.username, "TestUser"; got != want {
// 		t.Errorf("NewClient.username = %v, want %v", got, want)
// 	}
// 	if got, want := c.password, "TestPass"; got != want {
// 		t.Errorf("NewClient.password = %v, want %v", got, want)
// 	}

// 	// test passing http client as option
// 	if got, want := c.client, httpClient; got != want {
// 		t.Errorf("NewClient.client = %p, want %p", got, want)
// 	}
// }

func TestErrorResponse(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{
			Request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "foo.com",
				},
			},
			StatusCode: 500,
		},
		Message: "service gone",
		Err:     "boom",
	}
	if got, want := errResp.Error(), "GET http://foo.com: 500 service gone boom"; got != want {
		t.Errorf("ErrorResponse.Error() = %q, want %q", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody, outBody := func(v *url.Values) { v.Set("name", "bar") }, "name=bar"

	req, err := c.NewRequest("GET", inURL, inBody)
	if err != nil {
		t.Fatalf("NewRequest(%q) returned error: %v", inURL, err)
	}

	// test that the endpoint URL was correctly added to the base URL
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL = %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest("+`func(v *url.Values) { v.Set("name", "bar")`+") Body \nhave: %q\nwant: %q", got, want)
	}

	// test that Content-Type header is correctly set when body is set
	if got, want := req.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; got != want {
		t.Errorf("NewRequest() Content-Type \nhave: %q\nwant: %q", got, want)
	}
}

func TestClient_NewRequest_invalidMethod(t *testing.T) {
	s := strings.Split(runtime.Version(), ".")
	// This test requires Go version 1.7 or higher.
	if len(s) >= 2 && s[0] == "go1" && s[1] == "7" {
		c := NewClient(nil)
		_, err := c.NewRequest("invalid method", "/foo")
		if err == nil {
			t.Error("NewRequest with invalid method expected to return err")
		}
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	type foo struct {
		Bar string `json:"bar"`
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if want := "GET"; r.Method != want {
			t.Errorf("request method = %v, want %v", r.Method, want)
		}
		fmt.Fprint(w, `{"bar":"&bull; foobar"}`)
	})

	req, _ := client.NewRequest("GET", "/")

	body := new(foo)
	ctx := context.Background()
	_, err := client.Do(ctx, req, body)
	if err != nil {
		t.Fatalf("Do() returned err = %v", err)
	}

	want := &foo{"&bull; foobar"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Do() response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest("GET", "/")

	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

type errTransport struct{}

func (e errTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("connection refused")
}

func TestDo_returnsError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	client.client = &http.Client{
		Transport: &errTransport{},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	req, _ := client.NewRequest("GET", "/")
	ctx := context.Background()
	_, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected connection refused error.")
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".")
	ctx := context.Background()
	_, err := client.Do(ctx, req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestDo_bodyImplementsIOWriter(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "foo bar")
	})

	var body bytes.Buffer

	req, _ := client.NewRequest("GET", ".")
	ctx := context.Background()
	_, err := client.Do(ctx, req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
	if got, want := body.String(), "foo bar"; got != want {
		t.Errorf("Response Body is %q, want %q", got, want)
	}
}

func TestDo_decodeError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "this is not JSON")
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".")
	ctx := context.Background()
	_, err := client.Do(ctx, req, &body)
	if err == nil {
		t.Fatal("Expected JSON decode error.")
	}
}

func TestDo_contextCanceled(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := client.NewRequest("GET", ".")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Do(ctx, req, nil)
	if err == nil {
		t.Fatalf("Expected context canceled error.")
	}
}

func TestDo_nilContext(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".")
	_, err := client.Do(nil, req, nil)
	if err == nil {
		t.Errorf("Expected context must be non-nil error")
	}
}

// func TestClient_post_invalidID(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/animelist/update/", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")
// 		testID(t, r, "0")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testContentType(t, r, "application/x-www-form-urlencoded")
// 		// zeroEntry defined in anime_test.go
// 		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, 3))
// 		http.Error(w, "Invalid ID", http.StatusNotImplemented)
// 	})

// 	response, err := client.post("api/animelist/update/", 0, AnimeEntry{Status: OnHold}, true)

// 	if err == nil {
// 		t.Errorf("Anime.Update invalid ID should return err")
// 	}

// 	if response == nil {
// 		t.Errorf("Anime.Update invalid ID should return also return response")
// 	}
// }

// func TestClient_delete_invalidID(t *testing.T) {
// 	client, mux, teardown := setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/animelist/delete/", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "DELETE")
// 		testID(t, r, "0")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		http.Error(w, "Invalid ID", http.StatusNotImplemented)
// 	})

// 	response, err := client.delete("api/animelist/delete/", 0, true)

// 	if err == nil {
// 		t.Errorf("Anime.Delete invalid ID should return err")
// 	}

// 	if response == nil {
// 		t.Errorf("Anime.Delete invalid ID should return also return response")
// 	}
// }

func TestClient_NewRequest_badEndpoint(t *testing.T) {
	c := NewClient(nil)
	inURL := "%foo"
	_, err := c.NewRequest("GET", inURL)
	if err == nil {
		t.Errorf("NewRequest(%q) should return parse err", inURL)
	}
}

// func TestClient_NewRequest_xmlEncodeError(t *testing.T) {
// 	c := NewClient(nil)
// 	in := func() {} // xml.Marshal cannot encode a func
// 	_, err := c.NewRequest("GET", "/foo", in)
// 	if err == nil {
// 		t.Errorf("NewRequest receiving a function as body should return XML encode err")
// 	}
// }
