package mal_test

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

func ExampleMangaService_Details() {
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

func ExampleMangaService_DeleteMyListItem() {
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

	resp, err := c.Manga.DeleteMyListItem(ctx, 401)
	if err != nil {
		fmt.Printf("Manga.DeleteMyListItem error: %v", err)
		return
	}
	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}
