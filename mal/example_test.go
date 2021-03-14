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

// anime examples

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
			if err := serveStubFile(w, filename); err != nil {
				http.Error(w, fmt.Sprintf(`{"message": "%s", "error":"internal"}`, err), http.StatusInternalServerError)
			}
		}
	}

	mux.HandleFunc("/anime", serveStubHandler("animeList.json"))
	mux.HandleFunc("/anime/967", serveStubHandler("animeDetails.json"))

	return server
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

// func ExampleAnimeService_Update() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	resp, err := c.Anime.Update(9989, mal.AnimeEntry{Status: mal.Completed, Score: 9})
// 	if err != nil {
// 		log.Fatalf("Anime.Update error: %v, received: '%v'\n", err, string(resp.Body))
// 	}
// }

// func ExampleAnimeService_Search() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	result, resp, err := c.Anime.Search("anohana")
// 	if err != nil {
// 		log.Fatalf("Anime.Search error: %v, received: '%v'\n", err, string(resp.Body))
// 	}

// 	// For more complex searches, you can provide the % operator which is
// 	// escaped as %% in Go. Note: This is an undocumented API feature.
// 	//
// 	// As an example, if you search for "fate%%heaven%%flower" you can get one
// 	// accurate result of the title:
// 	//
// 	// Fate/stay night Movie: Heaven's Feel - I. presage flower

// 	// printing results
// 	for _, entry := range result.Rows {
// 		fmt.Printf("----------------------------------------\n")
// 		fmt.Printf("| ID: %v\n", entry.ID)
// 		fmt.Printf("| Title: %v\n", entry.Title)
// 		fmt.Printf("| Episodes: %d\n", entry.Episodes)
// 		fmt.Printf("| Type: %v\n", entry.Type)
// 		fmt.Printf("| Score: %v\n", entry.Score)
// 		fmt.Printf("| Synopsis: %v\n", entry.Synopsis)
// 		fmt.Printf("----------------------------------------\n")
// 		fmt.Printf("\n")
// 	}
// }

// func ExampleAnimeService_List() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	list, resp, err := c.Anime.List("Xinil")
// 	if err != nil {
// 		log.Fatalf("Anime.List error: %v, received: '%v'\n", err, string(resp.Body))
// 	}

// 	// printing results
// 	for _, anime := range list.Anime {
// 		fmt.Printf("----------------Anime-------------------\n")
// 		fmt.Printf("| MyID: %v\n", anime.MyID)
// 		fmt.Printf("| MyStartDate: %v\n", anime.MyStartDate)
// 		fmt.Printf("| MyFinishDate: %v\n", anime.MyFinishDate)
// 		fmt.Printf("| MyScore: %v\n", anime.MyScore)
// 		fmt.Printf("| MyStatus: %v\n", anime.MyStatus)
// 		fmt.Printf("| MyRewatching: %v\n", anime.MyRewatching)
// 		fmt.Printf("| MyRewatchingEp: %v\n", anime.MyRewatchingEp)
// 		fmt.Printf("| MyLastUpdated: %v\n", anime.MyLastUpdated)
// 		fmt.Printf("| MyTags: %v\n", anime.MyTags)
// 		fmt.Printf("| MyWatchedEpisodes: %v\n", anime.MyWatchedEpisodes)
// 		fmt.Printf("| SeriesAnimeDBID: %v\n", anime.SeriesAnimeDBID)
// 		fmt.Printf("| SeriesEpisodes: %v\n", anime.SeriesEpisodes)
// 		fmt.Printf("| SeriesTitle: %v\n", anime.SeriesTitle)
// 		fmt.Printf("| SeriesSynonyms: %v\n", anime.SeriesSynonyms)
// 		fmt.Printf("| SeriesType: %v\n", anime.SeriesType)
// 		fmt.Printf("| SeriesStatus: %v\n", anime.SeriesStatus)
// 		fmt.Printf("| SeriesStart: %v\n", anime.SeriesStart)
// 		fmt.Printf("| SeriesEnd: %v\n", anime.SeriesEnd)
// 		fmt.Printf("| SeriesImage: %v\n", anime.SeriesImage)
// 		fmt.Printf("\n")
// 	}
// 	fmt.Printf("----------------MyInfo------------------\n")
// 	fmt.Printf("| ID: %v\n", list.MyInfo.ID)
// 	fmt.Printf("| Name: %v\n", list.MyInfo.Name)
// 	fmt.Printf("| Completed: %v\n", list.MyInfo.Completed)
// 	fmt.Printf("| OnHold: %v\n", list.MyInfo.OnHold)
// 	fmt.Printf("| Dropped: %v\n", list.MyInfo.Dropped)
// 	fmt.Printf("| DaysSpentWatching: %v\n", list.MyInfo.DaysSpentWatching)
// 	fmt.Printf("| Watching: %v\n", list.MyInfo.Watching)
// 	fmt.Printf("| PlanToWatch: %v\n", list.MyInfo.PlanToWatch)
// 	fmt.Printf("----------------------------------------\n")
// }

// // manga examples

// func ExampleMangaService_Add() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	resp, err := c.Manga.Add(35733, mal.MangaEntry{Status: mal.Current, Chapter: 1, Volume: 1})
// 	if err != nil {
// 		log.Fatalf("Manga.Add error: %v, received: '%v'\n", err, string(resp.Body))
// 	}
// }

// func ExampleMangaService_Update() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	resp, err := c.Manga.Update(35733, mal.MangaEntry{Status: mal.Completed, Score: 9})
// 	if err != nil {
// 		log.Fatalf("Manga.Update error: %v, received: '%v'\n", err, string(resp.Body))
// 	}
// }

// func ExampleMangaService_Search() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	result, resp, err := c.Manga.Search("anohana")
// 	if err != nil {
// 		log.Fatalf("Manga.Search error: %v, received: '%v'\n", err, string(resp.Body))
// 	}

// 	for _, entry := range result.Rows {
// 		fmt.Printf("----------------------------------------\n")
// 		fmt.Printf("| ID: %v\n", entry.ID)
// 		fmt.Printf("| Title: %v\n", entry.Title)
// 		fmt.Printf("| Chapters: %d\n", entry.Chapters)
// 		fmt.Printf("| Volumes: %d\n", entry.Volumes)
// 		fmt.Printf("| Type: %v\n", entry.Type)
// 		fmt.Printf("| Score: %v\n", entry.Score)
// 		fmt.Printf("| Synopsis: %v\n", entry.Synopsis)
// 		fmt.Printf("----------------------------------------\n")
// 		fmt.Printf("\n")
// 	}
// }

// func ExampleMangaService_List() {
// 	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

// 	list, resp, err := c.Manga.List("Xinil")
// 	if err != nil {
// 		log.Fatalf("Manga.List error: %v, received: '%v'\n", err, string(resp.Body))
// 	}

// 	// printing results
// 	for _, manga := range list.Manga {
// 		fmt.Printf("----------------Manga-------------------\n")
// 		fmt.Printf("| MyID: %v\n", manga.MyID)
// 		fmt.Printf("| MyStartDate: %v\n", manga.MyStartDate)
// 		fmt.Printf("| MyFinishDate: %v\n", manga.MyFinishDate)
// 		fmt.Printf("| MyScore: %v\n", manga.MyScore)
// 		fmt.Printf("| MyStatus: %v\n", manga.MyStatus)
// 		fmt.Printf("| MyRereading: %v\n", manga.MyRereading)
// 		fmt.Printf("| MyRereadingChap: %v\n", manga.MyRereadingChap)
// 		fmt.Printf("| MyLastUpdated: %v\n", manga.MyLastUpdated)
// 		fmt.Printf("| MyTags: %v\n", manga.MyTags)
// 		fmt.Printf("| MyReadChapters: %v\n", manga.MyReadChapters)
// 		fmt.Printf("| MyReadVolumes: %v\n", manga.MyReadVolumes)
// 		fmt.Printf("| SeriesMangaDBID: %v\n", manga.SeriesMangaDBID)
// 		fmt.Printf("| SeriesChapters: %v\n", manga.SeriesChapters)
// 		fmt.Printf("| SeriesVolumes: %v\n", manga.SeriesVolumes)
// 		fmt.Printf("| SeriesTitle: %v\n", manga.SeriesTitle)
// 		fmt.Printf("| SeriesSynonyms: %v\n", manga.SeriesSynonyms)
// 		fmt.Printf("| SeriesType: %v\n", manga.SeriesType)
// 		fmt.Printf("| SeriesStatus: %v\n", manga.SeriesStatus)
// 		fmt.Printf("| SeriesStart: %v\n", manga.SeriesStart)
// 		fmt.Printf("| SeriesEnd: %v\n", manga.SeriesEnd)
// 		fmt.Printf("| SeriesImage: %v\n", manga.SeriesImage)
// 		fmt.Printf("\n")
// 	}
// 	fmt.Printf("----------------MyInfo------------------\n")
// 	fmt.Printf("| ID: %v\n", list.MyInfo.ID)
// 	fmt.Printf("| Name: %v\n", list.MyInfo.Name)
// 	fmt.Printf("| Completed: %v\n", list.MyInfo.Completed)
// 	fmt.Printf("| OnHold: %v\n", list.MyInfo.OnHold)
// 	fmt.Printf("| Dropped: %v\n", list.MyInfo.Dropped)
// 	fmt.Printf("| DaysSpentWatching: %v\n", list.MyInfo.DaysSpentWatching)
// 	fmt.Printf("| Reading: %v\n", list.MyInfo.Reading)
// 	fmt.Printf("| PlanToRead: %v\n", list.MyInfo.PlanToRead)
// 	fmt.Printf("----------------------------------------\n")
// }
