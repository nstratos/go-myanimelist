// +build integration

package mal_test

import (
	"flag"
	"testing"

	"github.com/nstratos/go-myanimelist/mal"
)

var (
	malUsername = flag.String("username", "testgopher", "MyAnimeList.net username to use for integration tests")
	malPassword = flag.String("password", "", "MyAnimeList.net password to use for integration tests")
	userAgent   = flag.String("agent", "", "User-Agent to use for integration tests")

	testAnimeIDs      = []int{1, 5, 6, 7}
	animeAddStatus    = mal.StatusPlanToWatch
	animeUpdateStatus = mal.StatusWatching

	testMangaIDs      = []int{1, 2, 3, 4}
	mangaAddStatus    = mal.StatusPlanToRead
	mangaUpdateStatus = mal.StatusReading

	// client is the MyAnimeList client being used for the integration tests.
	client *mal.Client
)

func init() {
	flag.Parse()
}

func setup(t *testing.T) {
	if *malPassword == "" {
		t.Errorf("No password provided for user %q.", *malUsername)
		t.Error("These tests are meant to be run with a dedicated test account.")
		t.Fatal("You might want to use: go test -tags=integration -username '<your account>' -password '<your password>'")
	}

	// Create mal client for tests.
	client = mal.NewClient(nil)
	client.SetCredentials(*malUsername, *malPassword)
	if *userAgent != "" {
		client.SetUserAgent(*userAgent)
	}
}

func TestAnimeServiceIntegration(t *testing.T) {
	setup(t)

	// Get anime list for test account. No authentication needed.
	list, _, err := client.Anime.List(*malUsername)
	if err != nil {
		t.Fatal("client.Anime.List returned err:", err)
	}

	// Being strict here. Anime list of test account must be empty.
	if len(list.Anime) != 0 {
		t.Fatalf("MyAnimeList.net test account %q is supposed to have 0 anime but has %d", *malUsername, len(list.Anime))
	}

	// Test if password matches test account username.
	if _, _, err := client.Account.Verify(); err != nil {
		t.Fatalf("Account verification for user %q failed: %v", *malUsername, err)
	}

	// Clean up all anime at the end.
	defer func() {
		for _, id := range testAnimeIDs {
			if _, derr := client.Anime.Delete(id); derr != nil {
				t.Errorf("client.Anime.Delete(%d) returned err: %v", id, derr)
			}
		}
	}()

	// Test adding all the anime.
	for _, id := range testAnimeIDs {
		if _, err := client.Anime.Add(id, mal.AnimeEntry{Status: animeAddStatus}); err != nil {
			t.Fatalf("client.Anime.Add(%d) returned err: %v", id, err)
		}
	}

	// Test updating all the anime.
	for _, id := range testAnimeIDs {
		if _, err := client.Anime.Update(id, mal.AnimeEntry{Status: animeUpdateStatus}); err != nil {
			t.Fatalf("client.Anime.Update(%d) returned err: %v", id, err)
		}
	}

	// Get anime list of test account for a second time.
	list, _, err = client.Anime.List(*malUsername)
	if err != nil {
		t.Fatal("client.Anime.List after additions returned err:", err)
	}

	// And make sure it has the number of anime it's supposed to have.
	if got, want := len(list.Anime), len(testAnimeIDs); got != want {
		t.Fatalf("Test account Anime number after additions = %d, want %d", got, want)
	}

	// And that they all have been updated appropriately.
	for _, anime := range list.Anime {
		if got, want := anime.MyStatus, animeUpdateStatus; got != want {
			t.Errorf("Anime ID: %d status = %d, want %d", anime.SeriesAnimeDBID, got, want)
		}
	}
}

func TestMangaServiceIntegration(t *testing.T) {
	setup(t)

	// Get manga list for test account. No authentication needed.
	list, _, err := client.Manga.List(*malUsername)
	if err != nil {
		t.Fatal("client.Manga.List returned err:", err)
	}

	// Being strict here. Manga list of test account must be empty.
	if len(list.Manga) != 0 {
		t.Fatalf("MyMangaList.net test account %q is supposed to have 0 manga but has %d", *malUsername, len(list.Manga))
	}

	// Test if password matches test account username.
	if _, _, err := client.Account.Verify(); err != nil {
		t.Fatalf("Account verification for user %q failed: %v", *malUsername, err)
	}

	// Clean up all manga at the end.
	defer func() {
		for _, id := range testMangaIDs {
			if _, derr := client.Manga.Delete(id); derr != nil {
				t.Errorf("client.Manga.Delete(%d) returned err: %v", id, derr)
			}
		}
	}()

	// Test adding all the manga.
	for _, id := range testMangaIDs {
		if _, err := client.Manga.Add(id, mal.MangaEntry{Status: mangaAddStatus}); err != nil {
			t.Fatalf("client.Manga.Add(%d) returned err: %v", id, err)
		}
	}

	// Test updating all the manga.
	for _, id := range testMangaIDs {
		if _, err := client.Manga.Update(id, mal.MangaEntry{Status: mangaUpdateStatus}); err != nil {
			t.Fatalf("client.Manga.Update(%d) returned err: %v", id, err)
		}
	}

	// Get manga list of test account for a second time.
	list, _, err = client.Manga.List(*malUsername)
	if err != nil {
		t.Fatal("client.Manga.List after additions returned err:", err)
	}

	// And make sure it has the number of manga it's supposed to have.
	if got, want := len(list.Manga), len(testMangaIDs); got != want {
		t.Fatalf("Test account Manga number after additions = %d, want %d", got, want)
	}

	// And that they all have been updated appropriately.
	for _, manga := range list.Manga {
		if got, want := manga.MyStatus, mangaUpdateStatus; got != want {
			t.Errorf("Manga ID: %d status = %d, want %d", manga.SeriesMangaDBID, got, want)
		}
	}
}
