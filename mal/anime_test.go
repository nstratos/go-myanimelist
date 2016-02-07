package mal

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var zeroEntry = "<entry><episode>0</episode><status>%v</status><score>0</score><times_rewatched>0</times_rewatched><enable_rewatching>0</enable_rewatching><comments></comments></entry>"

func TestAnimeService_Delete(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/delete/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testID(t, r, "55")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
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
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, "watching"))
		fmt.Fprintf(w, "Created")
	})

	_, err := client.Anime.Add(55, AnimeEntry{Status: "watching"})
	if err != nil {
		t.Errorf("Anime.Add returned error %v", err)
	}
}

func TestAnimeService_Update(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/animelist/update/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testID(t, r, "55")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testContentType(t, r, "application/x-www-form-urlencoded")
		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, "onhold"))
		fmt.Fprintf(w, "Updated")
	})

	_, err := client.Anime.Update(55, AnimeEntry{Status: "onhold"})
	if err != nil {
		t.Errorf("Anime.Update returned error %v", err)
	}
}

func TestAnimeService_Update_invalidID(t *testing.T) {
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
		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, "onhold"))
		http.Error(w, "Invalid ID", http.StatusNotImplemented)
	})

	response, err := client.Anime.Update(0, AnimeEntry{Status: "onhold"})

	if err == nil {
		t.Errorf("Anime.Update invalid ID should return err")
	}

	if response == nil {
		t.Errorf("Anime.Update invalid ID should return also return response")
	}
}

func TestAnimeService_Search(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/anime/search.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testURLValues(t, r, urlValues{"q": "query"})
		fmt.Fprintf(w, `
			<anime>
				<entry>
					<title>title1</title>
					<id>55</id>
				</entry>
				<entry>
					<title>title2</title>
					<id>56</id>
				</entry>
			</anime>`)
	})

	result, _, err := client.Anime.Search("query")
	if err != nil {
		t.Errorf("Anime.Search returned error %v", err)
	}
	want := &AnimeResult{
		[]AnimeRow{
			{ID: 55, Title: "title1"},
			{ID: 56, Title: "title2"},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Anime.Search returned %+v, want %+v", result, want)
	}
}

func TestAnimeService_Search_noContent(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/anime/search.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		testURLValues(t, r, urlValues{"q": "foo"})
		http.Error(w, "no content", http.StatusNoContent)
	})

	result, _, err := client.Anime.Search("foo")

	if err == nil {
		t.Errorf("Anime.Search for non existent query expected to return err")
	}

	if got, want := err, ErrNoContent; got != want {
		t.Errorf("Anime.Search for non existent query returned err %v, want %v", got, want)
	}

	if got := result; got != nil {
		t.Errorf("Anime.Search for non existent query returned result = %v, want %v", got, nil)
	}
}

func TestAnimeService_List(t *testing.T) {
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
			"type":   "anime",
			"u":      "AnotherTestUser",
		})
		fmt.Fprintf(w, `
			<myanimelist>
				<myinfo>
					<user_id>56</user_id>
					<user_name>AnotherTestUser</user_name>
				</myinfo>
				<anime>
					<series_animedb_id>1</series_animedb_id>
					<series_title>series title</series_title>
					<my_id>1234</my_id>
					<my_status>3</my_status>
				</anime>
			</myanimelist>
			`)
	})

	result, _, err := client.Anime.List("AnotherTestUser")
	if err != nil {
		t.Errorf("Anime.List returned error %v", err)
	}
	want := &AnimeList{
		MyInfo: AnimeMyInfo{ID: 56, Name: "AnotherTestUser"},
		Anime: []Anime{
			{
				SeriesAnimeDBID: 1,
				SeriesTitle:     "series title",
				MyID:            1234,
				MyStatus:        3,
			},
		},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Anime.List returned %+v, want %+v", result, want)
	}
}

func TestAnimeService_List_invalidUsername(t *testing.T) {
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
			"type":   "anime",
			"u":      "InvalidUser",
		})
		fmt.Fprintf(w, `
			<myanimelist>
				<error>Invalid username</error>
			</myanimelist>
			`)
	})

	result, _, err := client.Anime.List("InvalidUser")

	if err == nil {
		t.Errorf("Anime.List for invalid user expected to return err")
	}

	want := &AnimeList{Error: "Invalid username"}
	if !reflect.DeepEqual(want, result) {
		t.Errorf("Anime.List for invalid user returned result = %v, want %v", result, want)
	}
}

func TestAnimeService_List_httpError(t *testing.T) {
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
			"type":   "anime",
			"u":      "TestUser",
		})
		http.Error(w, "something broke", http.StatusInternalServerError)
	})

	result, _, err := client.Anime.List("TestUser")

	if err == nil {
		t.Errorf("Anime.List for server error expected to return err")
	}

	if got := result; got != nil {
		t.Errorf("Anime.List for server error returned result = %v, want %v", got, nil)
	}
}
