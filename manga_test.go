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
		testBasicAuth(t, r, true, "TestUser", "TestPass")
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
		testBasicAuth(t, r, true, "TestUser", "TestPass")
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
		testBasicAuth(t, r, true, "TestUser", "TestPass")
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
		testBasicAuth(t, r, true, "TestUser", "TestPass")
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

func TestMangaService_List(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/malappinfo.php", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testURLValues(t, r, urlValues{
			"status": "all",
			"type":   "manga",
			"u":      "AnotherTestUser",
		})
		fmt.Fprintf(w, `
			<mymangalist>
				<myinfo>
					<user_id>56</user_id>
					<user_name>AnotherTestUser</user_name>
				</myinfo>
				<manga>
					<series_mangadb_id>1</series_mangadb_id>
					<series_title>series title</series_title>
					<my_id>1234</my_id>
					<my_status>3</my_status>
				</manga>
			</mymangalist>
			`)
	})

	result, _, err := client.Manga.List("AnotherTestUser")
	if err != nil {
		t.Errorf("Manga.List returned error %v", err)
	}
	want := &MangaList{
		MyInfo: MyMangaInfo{ID: 56, Name: "AnotherTestUser"},
		Manga: []Manga{
			Manga{
				SeriesMangaDBID: 1,
				SeriesTitle:     "series title",
				MyID:            1234,
				MyStatus:        3,
			},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Manga.List returned %+v, want %+v", result, want)
	}
}
