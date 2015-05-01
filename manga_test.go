package mal

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMangaService_Delete(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/mangalist/delete/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testID(t, r, "55")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		fmt.Fprintf(w, "Deleted")
	})

	_, err := client.Manga.Delete(55)
	if err != nil {
		t.Errorf("Manga.Delete returned error %v", err)
	}
}

func TestMangaService_Add(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/mangalist/add/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testID(t, r, "55")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		testFormValue(t, r, "data", "<entry><status>watching</status></entry>")
		fmt.Fprintf(w, "Created")
	})

	_, err := client.Manga.Add(55, MangaEntry{Status: "watching"})
	if err != nil {
		t.Errorf("Manga.Add returned error %v", err)
	}
}

func TestMangaService_Update(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/mangalist/update/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testID(t, r, "55")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		testFormValue(t, r, "data", "<entry><status>onhold</status></entry>")
		fmt.Fprintf(w, "Updated")
	})

	_, err := client.Manga.Update(55, MangaEntry{Status: "onhold"})
	if err != nil {
		t.Errorf("Manga.Update returned error %v", err)
	}
}

func TestMangaService_Search(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/manga/search.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testURLValues(t, r, urlValues{"q": "query"})
		fmt.Fprintf(w, `
			<manga>
				<entry>
					<title>title1</title>
					<id>55</id>
				</entry>
				<entry>
					<title>title2</title>
					<id>56</id>
				</entry>
			</manga>`)
	})

	result, _, err := client.Manga.Search("query")
	if err != nil {
		t.Errorf("Manga.Search returned error %v", err)
	}
	want := &MangaResult{
		[]MangaRow{
			MangaRow{Row: Row{ID: 55, Title: "title1"}},
			MangaRow{Row: Row{ID: 56, Title: "title2"}},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Manga.Search returned %+v, want %+v", result, want)
	}
}
