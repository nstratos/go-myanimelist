package mal_test

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/nstratos/go-myanimelist/mal"
)

// newTestServer creates a test server which we can make it act as if it was
// the MyAnimeList API server. By contacting this test server instead of the
// MyAnimeList server we can write a runnable example which produces output.
// Unless you are planning to write tests or runnable examples, you can ignore
// this function's code.
func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/malappinfo.php", func(w http.ResponseWriter, r *http.Request) {
		list := mal.AnimeList{
			MyInfo: mal.AnimeMyInfo{Name: "testgopher"},
			Anime: []mal.Anime{
				{
					SeriesAnimeDBID: 1,
					SeriesTitle:     "anime title 1",
					MyStatus:        mal.Current,
				},
				{
					SeriesAnimeDBID: 2,
					SeriesTitle:     "anime title 2",
					MyStatus:        mal.Planned,
				},
				{
					SeriesAnimeDBID: 3,
					SeriesTitle:     "anime title 3",
					MyStatus:        mal.Planned,
				},
			},
		}
		err := xml.NewEncoder(w).Encode(list)
		if err != nil {
			log.Println(err)
		}
	})

	return server
}

func Example_testServer() {
	server := newTestServer()
	defer server.Close()

	// Here starts code that you would normally write when using this package.
	c := mal.NewClient(nil)
	// Except this line which makes the mal client contact our test server
	// instead of the MyAnimeList API.
	c.BaseURL, _ = url.Parse(server.URL)

	user := "testgopher"
	list, _, err := c.Anime.List(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s's anime list:\n", user)
	for _, anime := range list.Anime {
		switch anime.MyStatus {
		case mal.Current:
			fmt.Printf("Watching: %q\n", anime.SeriesTitle)
		case mal.Planned:
			fmt.Printf("Plan to watch: %q\n", anime.SeriesTitle)
		}
	}

	// Output:
	// testgopher's anime list:
	// Watching: "anime title 1"
	// Plan to watch: "anime title 2"
	// Plan to watch: "anime title 3"
}
