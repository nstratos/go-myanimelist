package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUserServiceAnimeList(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/foo/animelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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
		fmt.Fprintf(w, out)
	})

	ctx := context.Background()
	got, resp, err := client.User.AnimeList(ctx, "foo",
		AnimeStatusCompleted,
		SortAnimeListByAnimeID,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("User.AnimeList returned error: %v", err)
	}
	want := []AnimeWithStatus{
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
		testMethod(t, r, "GET")
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
