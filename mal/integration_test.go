package mal_test

import (
	"context"
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

var accessToken = flag.String("access-token", "", "MyAnimeList.net access token to use for integration tests")

func setup(ctx context.Context, t *testing.T) *mal.Client {
	if *accessToken == "" {
		t.Log("No access token provided.")
		t.Log("The integration tests are meant to be run with a dedicated test account with empty lists.")
		t.Skip("To run the integration tests use: go test --access-token '<your access token>'")
	}

	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	))
	return mal.NewClient(tc)
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	client := setup(ctx, t)

	username := testGetUserInfo(ctx, t, client)
	t.Run("UpdateUserAnimeList", func(t *testing.T) {
		testUpdateUserAnimeList(ctx, t, client, username)
	})
	t.Run("UpdateUserMangaList", func(t *testing.T) {
		testUpdateUserMangaList(ctx, t, client, username)
	})
}

func testGetUserInfo(ctx context.Context, t *testing.T, client *mal.Client) (username string) {
	t.Helper()
	// Get user info to find the username.
	info, _, err := client.User.MyInfo(ctx)
	if err != nil {
		t.Fatal("User.MyInfo returned err:", err)
	}

	username = info.Name
	t.Logf("Running integration tests using user: %q", username)
	return username
}

func testUpdateUserAnimeList(ctx context.Context, t *testing.T, client *mal.Client, username string) {
	// Get anime list for test account.
	const me = "@me"
	list, _, err := client.User.AnimeList(ctx, me)
	if err != nil {
		t.Fatalf("User.AnimeList(%q) returned err: %s", me, err)
	}

	// Being strict here. Anime list of test account must be empty.
	if len(list) != 0 {
		t.Fatalf("MyAnimeList.net test account %q is supposed to have 0 anime but has %d", username, len(list))
	}

	testAnimeIDs := []int{1, 5, 6, 7}
	// Clean up all anime at the end.
	defer func() {
		for _, id := range testAnimeIDs {
			if _, delErr := client.Anime.DeleteMyListItem(ctx, id); delErr != nil {
				t.Errorf("Anime.DeleteMyListItem(%d) returned err: %v", id, delErr)
			}
		}
	}()

	// Test adding some anime.
	for _, id := range testAnimeIDs {
		if _, _, err := client.Anime.UpdateMyListStatus(ctx, id,
			mal.AnimeStatusWatching,
			mal.Comments("test comment"),
			mal.IsRewatching(true),
			mal.NumEpisodesWatched(1),
			mal.NumTimesRewatched(1),
			mal.Priority(1),
			mal.RewatchValue(1),
			mal.Score(1),
			mal.Tags{"foo", "bar"},
		); err != nil {
			t.Fatalf("Anime.UpdateMyListStatus(%d) returned err: %v", id, err)
		}
	}

	// Get anime list of test account for a second time.
	list, _, err = client.User.AnimeList(ctx, me,
		mal.Fields{"list_status{num_times_rewatched, rewatch_value, priority, comments, tags}"},
	)
	if err != nil {
		t.Fatalf("User.AnimeList(%q) after additions returned err: %s", me, err)
	}

	// And make sure it has the number of anime it's supposed to have.
	if got, want := len(list), len(testAnimeIDs); got != want {
		t.Fatalf("Test account Anime number after additions = %d, want %d", got, want)
	}

	// And that they all have been updated appropriately.
	for _, a := range list {
		want := mal.AnimeListStatus{
			Status:             mal.AnimeStatusWatching,
			Score:              1,
			NumEpisodesWatched: 1,
			IsRewatching:       true,
			Priority:           1,
			NumTimesRewatched:  1,
			RewatchValue:       1,
			Tags:               []string{"foo", "bar"},
			Comments:           "test comment",
		}
		a.Status.UpdatedAt = time.Time{}
		if got := a.Status; !reflect.DeepEqual(got, want) {
			t.Errorf("Anime ID: %d AnimeListStatus\nhave: %+v\nwant: %+v", a.Anime.ID, got, want)
		}
	}
}

func testUpdateUserMangaList(ctx context.Context, t *testing.T, client *mal.Client, username string) {
	// Get manga list for test account.
	const me = "@me"
	list, _, err := client.User.MangaList(ctx, me)
	if err != nil {
		t.Fatalf("User.MangaList(%q) returned err: %s", me, err)
	}

	// Being strict here. Manga list of test account must be empty.
	if len(list) != 0 {
		t.Fatalf("MyMangaList.net test account %q is supposed to have 0 manga but has %d", username, len(list))
	}

	testMangaIDs := []int{1, 2, 3, 4}
	// Clean up all manga at the end.
	defer func() {
		for _, id := range testMangaIDs {
			if _, delErr := client.Manga.DeleteMyListItem(ctx, id); delErr != nil {
				t.Errorf("Manga.DeleteMyListItem(%d) returned err: %v", id, delErr)
			}
		}
	}()

	// Test adding some manga.
	for _, id := range testMangaIDs {
		if _, _, err := client.Manga.UpdateMyListStatus(ctx, id,
			mal.MangaStatusReading,
			mal.Comments("test comment"),
			mal.IsRereading(true),
			mal.NumChaptersRead(1),
			mal.NumVolumesRead(1),
			mal.NumTimesReread(1),
			mal.Priority(1),
			mal.RereadValue(1),
			mal.Score(1),
			mal.Tags{"foo", "bar"},
		); err != nil {
			t.Fatalf("Manga.UpdateMyListStatus(%d) returned err: %v", id, err)
		}
	}

	// Get manga list of test account for a second time.
	list, _, err = client.User.MangaList(ctx, me,
		mal.Fields{"list_status{num_times_reread, reread_value, priority, comments, tags}"},
	)
	if err != nil {
		t.Fatalf("User.MangaList(%q) after additions returned err: %s", me, err)
	}

	// And make sure it has the number of manga it's supposed to have.
	if got, want := len(list), len(testMangaIDs); got != want {
		t.Fatalf("Test account Manga number after additions = %d, want %d", got, want)
	}

	// And that they all have been updated appropriately.
	for _, a := range list {
		want := mal.MangaListStatus{
			Status:          mal.MangaStatusReading,
			Score:           1,
			NumChaptersRead: 1,
			NumVolumesRead:  1,
			IsRereading:     true,
			Priority:        1,
			NumTimesReread:  1,
			RereadValue:     1,
			Tags:            []string{"foo", "bar"},
			Comments:        "test comment",
		}
		a.Status.UpdatedAt = time.Time{}
		if got := a.Status; !reflect.DeepEqual(got, want) {
			t.Errorf("Manga ID: %d MangaListStatus\nhave: %+v\nwant: %+v", a.Manga.ID, got, want)
		}
	}
}

func testForumMethods(ctx context.Context, t *testing.T, client *mal.Client, username string) {
	_, _, err := client.Forum.Boards(ctx)
	if err != nil {
		t.Errorf("Forum.Boards returned error: %v", err)
	}

	topics, _, err := client.Forum.Topics(ctx, mal.Query("kiseijuu"), mal.SortTopicsRecent)
	if err != nil {
		t.Errorf("Forum.Topics returned error: %v", err)
	}
	if len(topics) == 0 {
		t.Fatal("Forum.Topics returned 0 topics")
	}

	_, _, err = client.Forum.TopicDetails(ctx, topics[0].ID)
	if err != nil {
		t.Errorf("Forum.TopicDetails returned error: %v", err)
	}
}
