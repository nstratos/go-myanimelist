package mal_test

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
)

func ExampleMangaService_List() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	manga, _, err := c.Manga.List(ctx, "parasyte",
		mal.Fields{"num_volumes", "num_chapters", "alternative_titles"},
		mal.Limit(3),
		mal.Offset(0),
	)
	if err != nil {
		fmt.Printf("Manga.List error: %v", err)
		return
	}
	for _, m := range manga {
		fmt.Printf("ID: %5d, Volumes: %3d, Chapters: %3d %s (%s)\n", m.ID, m.NumVolumes, m.NumChapters, m.Title, m.AlternativeTitles.En)
	}
	// Output:
	// ID:   401, Volumes:  10, Chapters:  64 Kiseijuu (Parasyte)
	// ID: 78789, Volumes:   1, Chapters:  12 Neo Kiseijuu (Neo Parasyte m)
	// ID: 80797, Volumes:   1, Chapters:  15 Neo Kiseijuu f (Neo Parasyte f)
}

func ExampleMangaService_Details() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	m, _, err := c.Manga.Details(ctx, 401,
		mal.Fields{
			"alternative_titles",
			"media_type",
			"num_volumes",
			"num_chapters",
			"authors{last_name, first_name}",
			"genres",
			"status",
		},
	)
	if err != nil {
		fmt.Printf("Manga.Details error: %v", err)
		return
	}

	fmt.Printf("%s\n", m.Title)
	fmt.Printf("ID: %d\n", m.ID)
	fmt.Printf("English: %s\n", m.AlternativeTitles.En)
	fmt.Printf("Type: %s\n", strings.Title(m.MediaType))
	fmt.Printf("Volumes: %d\n", m.NumVolumes)
	fmt.Printf("Chapters: %d\n", m.NumChapters)
	fmt.Print("Studios: ")
	delim := ""
	for _, s := range m.Authors {
		fmt.Printf("%s%s, %s (%s)", delim, s.Person.LastName, s.Person.FirstName, s.Role)
		delim = " "
	}
	fmt.Println()
	fmt.Print("Genres: ")
	delim = ""
	for _, g := range m.Genres {
		fmt.Printf("%s%s", delim, g.Name)
		delim = " "
	}
	fmt.Println()
	fmt.Printf("Status: %s\n", strings.Title(m.Status))
	// Output:
	// Kiseijuu
	// ID: 401
	// English: Parasyte
	// Type: Manga
	// Volumes: 10
	// Chapters: 64
	// Studios: Iwaaki, Hitoshi (Story & Art)
	// Genres: Action Psychological Sci-Fi Drama Horror Seinen
	// Status: Finished
}

func ExampleMangaService_Ranking() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	manga, _, err := c.Manga.Ranking(ctx,
		mal.MangaRankingByPopularity,
		mal.Fields{"rank", "popularity"},
		mal.Limit(6),
	)
	if err != nil {
		fmt.Printf("Manga.Ranking error: %v", err)
		return
	}
	for _, m := range manga {
		fmt.Printf("Rank: %5d, Popularity: %5d %s\n", m.Rank, m.Popularity, m.Title)
	}
	// Output:
	// Rank:    38, Popularity:     1 Shingeki no Kyojin
	// Rank:     3, Popularity:     2 One Piece
	// Rank:     1, Popularity:     3 Berserk
	// Rank:   566, Popularity:     4 Naruto
	// Rank:   106, Popularity:     5 Tokyo Ghoul
	// Rank:    39, Popularity:     6 One Punch-Man
}

func ExampleMangaService_DeleteMyListItem() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	resp, err := c.Manga.DeleteMyListItem(ctx, 401)
	if err != nil {
		fmt.Printf("Manga.DeleteMyListItem error: %v", err)
		return
	}
	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}
