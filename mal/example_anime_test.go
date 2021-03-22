package mal_test

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

//go:embed testdata/*.json
var testDataJSON embed.FS

// newStubServer creates a stub server which serves some premade responses. By
// contacting this server instead of the real API we can have runnable examples
// which always produce the same output.
func newStubServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	serveStubFile := func(w io.Writer, filename string) error {
		stubResponses, err := fs.Sub(testDataJSON, "testdata")
		if err != nil {
			return err
		}
		f, err := stubResponses.Open(filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
		return nil
	}

	serveStubHandler := func(filename string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			malError := func(err string) string {
				return fmt.Sprintf(`{"message": "", "error":"%s"}`, err)
			}
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusOK)
			case http.MethodGet, http.MethodPatch:
				if err := serveStubFile(w, filename); err != nil {
					http.Error(w, malError("internal"), http.StatusInternalServerError)
				}
			default:
				http.Error(w, malError("not_allowed"), http.StatusMethodNotAllowed)
			}
		}
	}

	mux.HandleFunc("/anime", serveStubHandler("animeList.json"))
	mux.HandleFunc("/anime/967", serveStubHandler("animeDetails.json"))
	mux.HandleFunc("/anime/967/my_list_status", serveStubHandler("updateMyAnimeList.json"))
	mux.HandleFunc("/anime/ranking", serveStubHandler("animeRanking.json"))
	mux.HandleFunc("/anime/season/2020/fall", serveStubHandler("animeSeasonal.json"))
	mux.HandleFunc("/anime/suggestions", serveStubHandler("animeSuggested.json"))
	mux.HandleFunc("/manga", serveStubHandler("mangaList.json"))
	mux.HandleFunc("/manga/401", serveStubHandler("mangaDetails.json"))
	mux.HandleFunc("/manga/401/my_list_status", serveStubHandler("updateMyMangaList.json"))
	mux.HandleFunc("/manga/ranking", serveStubHandler("mangaRanking.json"))
	mux.HandleFunc("/users/@me", serveStubHandler("userMyInfo.json"))
	mux.HandleFunc("/users/@me/animelist", serveStubHandler("userAnimeList.json"))
	mux.HandleFunc("/users/@me/mangalist", serveStubHandler("userMangaList.json"))
	mux.HandleFunc("/forum/topics", serveStubHandler("forumTopics.json"))
	mux.HandleFunc("/forum/topic/1877721", serveStubHandler("forumTopicDetails.json"))

	return server
}

func ExampleAnimeService_List() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	anime, _, err := c.Anime.List(ctx, "hokuto no ken",
		mal.Fields{"rank", "popularity", "start_season"},
		mal.Limit(5),
		mal.Offset(0),
	)
	if err != nil {
		fmt.Printf("Anime.List error: %v", err)
		return
	}
	for _, a := range anime {
		fmt.Printf("ID: %5d, Rank: %5d, Popularity: %5d %s (%d)\n", a.ID, a.Rank, a.Popularity, a.Title, a.StartSeason.Year)
	}
	// Output:
	// ID:   967, Rank:   556, Popularity:  1410 Hokuto no Ken (1984)
	// ID:  1356, Rank:  1423, Popularity:  3367 Hokuto no Ken 2 (1987)
	// ID:  1358, Rank:  2757, Popularity:  3964 Hokuto no Ken Movie (1986)
}

func ExampleAnimeService_Details() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	a, _, err := c.Anime.Details(ctx, 967,
		mal.Fields{
			"alternative_titles",
			"media_type",
			"num_episodes",
			"start_season",
			"source",
			"genres",
			"studios",
			"average_episode_duration",
		},
	)
	if err != nil {
		fmt.Printf("Anime.Details error: %v", err)
		return
	}

	fmt.Printf("%s\n", a.Title)
	fmt.Printf("ID: %d\n", a.ID)
	fmt.Printf("English: %s\n", a.AlternativeTitles.En)
	fmt.Printf("Type: %s\n", strings.ToUpper(a.MediaType))
	fmt.Printf("Episodes: %d\n", a.NumEpisodes)
	fmt.Printf("Premiered: %s %d\n", strings.Title(a.StartSeason.Season), a.StartSeason.Year)
	fmt.Print("Studios: ")
	delim := ""
	for _, s := range a.Studios {
		fmt.Printf("%s%s", delim, s.Name)
		delim = " "
	}
	fmt.Println()
	fmt.Printf("Source: %s\n", strings.Title(a.Source))
	fmt.Print("Genres: ")
	delim = ""
	for _, g := range a.Genres {
		fmt.Printf("%s%s", delim, g.Name)
		delim = " "
	}
	fmt.Println()
	fmt.Printf("Duration: %d min. per ep.\n", a.AverageEpisodeDuration/60)
	// Output:
	// Hokuto no Ken
	// ID: 967
	// English: Fist of the North Star
	// Type: TV
	// Episodes: 109
	// Premiered: Fall 1984
	// Studios: Toei Animation
	// Source: Manga
	// Genres: Action Drama Martial Arts Sci-Fi Shounen
	// Duration: 25 min. per ep.
}

func ExampleAnimeService_Ranking() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	anime, _, err := c.Anime.Ranking(ctx,
		mal.AnimeRankingAiring,
		mal.Fields{"rank", "popularity"},
		mal.Limit(6),
	)
	if err != nil {
		fmt.Printf("Anime.Ranking error: %v", err)
		return
	}
	for _, a := range anime {
		fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
	}
	// Output:
	// Rank:     2, Popularity:   104 Shingeki no Kyojin: The Final Season
	// Rank:    59, Popularity:   371 Re:Zero kara Hajimeru Isekai Seikatsu 2nd Season Part 2
	// Rank:    67, Popularity:  1329 Yuru Camp△ Season 2
	// Rank:    69, Popularity:   109 Jujutsu Kaisen (TV)
	// Rank:    83, Popularity:  3714 Holo no Graffiti
	// Rank:    85, Popularity:   304 Horimiya
}

func ExampleAnimeService_Seasonal() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	anime, _, err := c.Anime.Seasonal(ctx, 2020, mal.AnimeSeasonFall,
		mal.Fields{"rank", "popularity"},
		mal.SortSeasonalByAnimeNumListUsers,
		mal.Limit(3),
		mal.Offset(0),
	)
	if err != nil {
		fmt.Printf("Anime.Seasonal error: %v", err)
		return
	}
	for _, a := range anime {
		fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
	}
	// Output:
	// Rank:    93, Popularity:    31 One Piece
	// Rank:  1870, Popularity:    92 Black Clover
	// Rank:    62, Popularity:   106 Jujutsu Kaisen (TV)
}

func ExampleAnimeService_Suggested() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	anime, _, err := c.Anime.Suggested(ctx,
		mal.Limit(10),
		mal.Fields{"rank", "popularity"},
	)
	if err != nil {
		fmt.Printf("Anime.Suggested error: %v", err)
		return
	}
	for _, a := range anime {
		fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
	}
	// Output:
	// Rank:   971, Popularity:   916 Kill la Kill Specials
	// Rank:   726, Popularity:  4972 Osomatsu-san Movie
	// Rank:   943, Popularity:  4595 Watashi no Ashinaga Ojisan
}

func ExampleAnimeService_DeleteMyListItem() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	resp, err := c.Anime.DeleteMyListItem(ctx, 967)
	if err != nil {
		fmt.Printf("Anime.DeleteMyListItem error: %v", err)
		return
	}
	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}
