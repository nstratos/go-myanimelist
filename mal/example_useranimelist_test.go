package mal_test

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
)

func ExampleUserService_AnimeList() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	anime, _, err := c.User.AnimeList(ctx, "@me",
		mal.Fields{"list_status"},
		mal.SortAnimeListByListUpdatedAt,
		mal.Limit(5),
	)
	if err != nil {
		fmt.Printf("User.AnimeList error: %v", err)
		return
	}
	for _, a := range anime {
		fmt.Printf("ID: %5d, Status: %15q, Episodes Watched: %3d %s\n", a.Anime.ID, a.Status.Status, a.Status.NumEpisodesWatched, a.Anime.Title)
	}
	// Output:
	// ID:   967, Status:      "watching", Episodes Watched:  73 Hokuto no Ken
	// ID:   820, Status:      "watching", Episodes Watched:   1 Ginga Eiyuu Densetsu
	// ID: 42897, Status:      "watching", Episodes Watched:   2 Horimiya
	// ID:  1453, Status:      "watching", Episodes Watched:  28 Maison Ikkoku
	// ID: 37521, Status:     "completed", Episodes Watched:  24 Vinland Saga
}

func ExampleAnimeService_UpdateMyListStatus() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	s, _, err := c.Anime.UpdateMyListStatus(ctx, 967,
		mal.AnimeStatusWatching,
		mal.NumEpisodesWatched(73),
		mal.Score(8),
		mal.Comments("You wa shock!"),
		mal.StartDate(time.Date(2022, 02, 20, 0, 0, 0, 0, time.UTC)),
		mal.FinishDate(time.Time{}), // Remove an existing date.
	)
	if err != nil {
		fmt.Printf("Anime.UpdateMyListStatus error: %v", err)
		return
	}
	fmt.Printf("Status: %q, Score: %d, Episodes Watched: %d, Comments: %q, Start Date: %s\n", s.Status, s.Score, s.NumEpisodesWatched, s.Comments, s.StartDate)
	// Output:
	// Status: "watching", Score: 8, Episodes Watched: 73, Comments: "You wa shock!", Start Date: 2022-02-20
}
