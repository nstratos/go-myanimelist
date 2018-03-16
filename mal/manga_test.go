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
		testFormValue(t, r, "data", "<entry><status>1</status></entry>")
		fmt.Fprintf(w, "Created")
	})

	_, err := client.Manga.Add(55, MangaEntry{Status: Current})
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
		testFormValue(t, r, "data", "<entry><status>3</status></entry>")
		fmt.Fprintf(w, "Updated")
	})

	_, err := client.Manga.Update(55, MangaEntry{Status: OnHold})
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
			{ID: 55, Title: "title1"},
			{ID: 56, Title: "title2"},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Manga.Search returned %+v, want %+v", result, want)
	}
}

func TestMangaService_Search_noContent(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/manga/search.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testURLValues(t, r, urlValues{"q": "foo"})
		http.Error(w, "no content", http.StatusNoContent)
	})

	result, _, err := client.Manga.Search("foo")

	if err == nil {
		t.Errorf("Manga.Search for non existent query expected to return err")
	}

	if got, want := err, ErrNoContent; got != want {
		t.Errorf("Manga.Search for non existent query returned err %v, want %v", got, want)
	}

	if got := result; got != nil {
		t.Errorf("Manga.Search for non existent query returned result = %v, want %v", got, nil)
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
					<my_rereadingg>1</my_rereadingg>
					<my_rereading_chap>2</my_rereading_chap>
				</manga>
			</mymangalist>
			`)
	})

	got, _, err := client.Manga.List("AnotherTestUser")
	if err != nil {
		t.Errorf("Manga.List returned error %v", err)
	}
	want := &MangaList{
		MyInfo: MangaMyInfo{ID: 56, Name: "AnotherTestUser"},
		Manga: []Manga{
			{
				SeriesMangaDBID: 1,
				SeriesTitle:     "series title",
				MyID:            1234,
				MyStatus:        3,
				MyRereading:     1,
				MyRereadingChap: 2,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Manga.List \nhave: %#v\nwant: %#v", got, want)
	}
}

func TestMangaService_List_invalidUsername(t *testing.T) {
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
			"u":      "InvalidUser",
		})
		fmt.Fprintf(w, `
			<myanimelist>
				<error>Invalid username</error>
			</myanimelist>
			`)
	})

	result, _, err := client.Manga.List("InvalidUser")

	if err == nil {
		t.Errorf("Manga.List for invalid user expected to return err")
	}

	want := &MangaList{Error: "Invalid username"}
	if !reflect.DeepEqual(want, result) {
		t.Errorf("Manga.List for invalid user returned result = %v, want %v", result, want)
	}
}

func TestMangaService_List_httpError(t *testing.T) {
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
			"u":      "TestUser",
		})
		http.Error(w, "something broke", http.StatusInternalServerError)
	})

	result, _, err := client.Manga.List("TestUser")

	if err == nil {
		t.Errorf("Manga.List for server error expected to return err")
	}

	if got := result; got != nil {
		t.Errorf("Manga.List for server error returned result = %v, want %v", got, nil)
	}
}
