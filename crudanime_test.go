package mal

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestAnimeService_Delete(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/delete/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testID(t, r, "55")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		fmt.Fprintf(w, "Deleted")
	})

	_, err := client.Anime.Delete(55)
	if err != nil {
		t.Errorf("Anime.Delete returned error %v", err)
	}
}

func TestAnimeService_Add(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/add/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testID(t, r, "55")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		testFormValue(t, r, "data", "<entry><status>watching</status></entry>")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Created")
	})

	resp, err := client.Anime.Add(55, AnimeData{Status: "watching"})
	if err != nil {
		t.Errorf("Anime.Delete returned error %v", err)
	}
	testResponseStatus(t, resp, 201)
}

func testResponseStatus(t *testing.T, resp *http.Response, want int) {
	if resp.StatusCode != want {
		t.Errorf("status code = %d, want %d", resp.StatusCode, want)
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
