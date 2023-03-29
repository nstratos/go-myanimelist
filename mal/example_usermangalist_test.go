package mal_test

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
)

func ExampleUserService_MangaList() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	manga, _, err := c.User.MangaList(ctx, "@me",
		mal.Fields{"list_status"},
		mal.SortMangaListByListUpdatedAt,
		mal.Limit(2),
	)
	if err != nil {
		fmt.Printf("User.MangaList error: %v", err)
		return
	}
	for _, m := range manga {
		fmt.Printf("ID: %5d, Status: %15q, Volumes Read: %3d, Chapters Read: %3d %s\n", m.Manga.ID, m.Status.Status, m.Status.NumVolumesRead, m.Status.NumChaptersRead, m.Manga.Title)
	}
	// Output:
	// ID:    21, Status:     "completed", Volumes Read:  12, Chapters Read: 108 Death Note
	// ID:   401, Status:       "reading", Volumes Read:   1, Chapters Read:   5 Kiseijuu
}

func ExampleMangaService_UpdateMyListStatus() {
	ctx := context.Background()

	c := mal.NewClient(nil)

	// Ignore the 3 following lines. A stub server is used instead of the real
	// API to produce testable examples. See: https://go.dev/blog/examples
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	s, _, err := c.Manga.UpdateMyListStatus(ctx, 401,
		mal.MangaStatusReading,
		mal.NumVolumesRead(1),
		mal.NumChaptersRead(5),
		mal.Comments("Migi"),
		mal.StartDate(time.Date(2022, 02, 20, 0, 0, 0, 0, time.UTC)),
		mal.FinishDate(time.Time{}), // Remove an existing date.
	)
	if err != nil {
		fmt.Printf("Manga.UpdateMyListStatus error: %v", err)
		return
	}
	fmt.Printf("Status: %q, Volumes Read: %d, Chapters Read: %d, Comments: %q, Start Date: %s\n", s.Status, s.NumVolumesRead, s.NumChaptersRead, s.Comments, s.StartDate)
	// Output:
	// Status: "reading", Volumes Read: 1, Chapters Read: 5, Comments: "Migi", Start Date: 2022-02-20
}
