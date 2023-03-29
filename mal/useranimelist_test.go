package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestUserServiceAnimeList(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/foo/animelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"status": "completed",
			"sort":   "anime_id",
			"fields": "foo,bar",
			"limit":  "10",
			"offset": "0",
			"nsfw":   "true",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 },
			  "list_status": {
			    "status": "plan_to_watch"
			  }
		    },
		    {
		      "node": { "id": 2 },
			  "list_status": {
			    "status": "watching"
			  }
		    }
		  ],
		  "paging": {
		    "next": "?offset=4",
		    "previous": "?offset=2"
		  }
		}`
		fmt.Fprint(w, out)
	})

	ctx := context.Background()
	got, resp, err := client.User.AnimeList(ctx, "foo",
		AnimeStatusCompleted,
		SortAnimeListByAnimeID,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
		NSFW(true),
	)
	if err != nil {
		t.Errorf("User.AnimeList returned error: %v", err)
	}
	want := []UserAnime{
		{
			Anime:  Anime{ID: 1},
			Status: AnimeListStatus{Status: "plan_to_watch"},
		},
		{
			Anime:  Anime{ID: 2},
			Status: AnimeListStatus{Status: "watching"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.AnimeList returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 2, "User.AnimeList")
}

func TestUserServiceAnimeListError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/foo/animelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"mal is down","error":"internal"}`, 500)
	})

	ctx := context.Background()
	_, resp, err := client.User.AnimeList(ctx, "foo",
		AnimeStatusCompleted,
		SortAnimeListByAnimeID,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err == nil {
		t.Fatal("User.AnimeList expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusInternalServerError, "User.AnimeList")
	testErrorResponse(t, err, ErrorResponse{Message: "mal is down", Err: "internal"})
}

func TestAnimeServiceUpdateMyListStatus(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		testContentType(t, r, "application/x-www-form-urlencoded")
		testBody(t, r, "comments=comments&finish_date=&is_rewatching=true&num_times_rewatched=2&num_watched_episodes=3&priority=2&rewatch_value=1&score=8&start_date=2022-02-20&status=completed&tags=foo%2Cbar")
		const out = `
		{
		  "status": "completed",
		  "score": 8,
		  "num_episodes_watched": 3,
		  "is_rewatching": true,
		  "updated_at": "2018-04-25T15:59:52Z",
		  "start_date": "2022-02-20",
		  "priority": 2,
		  "num_times_rewatched": 2,
		  "rewatch_value": 1,
		  "tags": ["foo","bar"],
		  "comments": "comments"
		}`
		fmt.Fprint(w, out)
	})

	ctx := context.Background()
	got, _, err := client.Anime.UpdateMyListStatus(ctx, 1,
		AnimeStatusCompleted,
		IsRewatching(true),
		Score(8),
		NumEpisodesWatched(3),
		Priority(2),
		NumTimesRewatched(2),
		RewatchValue(1),
		Tags{"foo", "bar"},
		Comments("comments"),
		StartDate(time.Date(2022, 02, 20, 0, 0, 0, 0, time.UTC)),
		FinishDate(time.Time{}),
	)
	if err != nil {
		t.Errorf("Anime.UpdateMyListStatus returned error: %v", err)
	}

	want := &AnimeListStatus{
		Status:             AnimeStatusCompleted,
		IsRewatching:       true,
		Score:              8,
		NumEpisodesWatched: 3,
		Priority:           2,
		NumTimesRewatched:  2,
		RewatchValue:       1,
		Tags:               []string{"foo", "bar"},
		Comments:           "comments",
		UpdatedAt:          time.Date(2018, 04, 25, 15, 59, 52, 0, time.UTC),
		StartDate:          "2022-02-20",
		FinishDate:         "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.UpdateMyListStatus returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestAnimeServiceUpdateMyListStatusError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		http.Error(w, `{"message":"mal is down","error":"internal"}`, 500)
	})

	ctx := context.Background()
	_, resp, err := client.Anime.UpdateMyListStatus(ctx, 1)
	if err == nil {
		t.Fatal("Anime.UpdateMyListStatus expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusInternalServerError, "Anime.UpdateMyListStatus")
	testErrorResponse(t, err, ErrorResponse{Message: "mal is down", Err: "internal"})
}

func TestAnimeServiceDeleteMyListItem(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	ctx := context.Background()
	resp, err := client.Anime.DeleteMyListItem(ctx, 1)
	if err != nil {
		t.Errorf("Anime.DeleteMyListItem returned error: %v", err)
	}
	testResponseStatusCode(t, resp, http.StatusOK, "Anime.DeleteMyListItem")
}

func TestAnimeServiceDeleteMyListItemError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		http.Error(w, `{"message":"anime not found","error":"not_found"}`, http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Anime.DeleteMyListItem(ctx, 1)
	if err == nil {
		t.Fatal("Anime.DeleteMyListItem expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusNotFound, "Anime.DeleteMyListItem")
	testErrorResponse(t, err, ErrorResponse{Message: "anime not found", Err: "not_found"})
}
