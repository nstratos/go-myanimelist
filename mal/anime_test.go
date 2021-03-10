package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestAnimeServiceDetails(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"fields": "foo,bar",
		})
		testBody(t, r, "")
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	a, _, err := client.Anime.Details(ctx, 1, Fields{"foo,bar"})
	if err != nil {
		t.Errorf("Anime.Details returned error: %v", err)
	}
	want := &Anime{ID: 1}
	if got := a; !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Details returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestAnimeServiceDetailsError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"anime deleted","error":"not_found"}`, 404)
	})

	ctx := context.Background()
	_, _, err := client.Anime.Details(ctx, 1)
	if err == nil {
		t.Fatal("Anime.Details expected not found error, got no error.")
	}
	testErrorResponse(t, err, ErrorResponse{Message: "anime deleted", Err: "not_found"})
}

func TestAnimeServiceList(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"q":      "query",
			"fields": "foo,bar",
			"limit":  "10",
			"offset": "0",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 }
		    },
		    {
		      "node": { "id": 2 }
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
	got, resp, err := client.Anime.List(ctx, "query",
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("Anime.List returned error: %v", err)
	}
	want := []Anime{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.List returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 2, "Anime.List")
}

func TestAnimeServiceListParsePagingError(t *testing.T) {
	tests := []struct {
		name string
		out  string
	}{
		{
			name: "cannot parse next url",
			out: `{
			  "data": [],
			  "paging": { "next": "\f", "previous": "?offset=2" }
			}`,
		},
		{
			name: "cannot parse previous url",
			out: `{
			  "data": [],
			  "paging": { "next": "?offset=2", "previous": "\f" }
			}`,
		},
		{
			name: "cannot parse next offset as int",
			out: `{
			  "data": [],
			  "paging": { "next": "?offset=foo", "previous": "?offset=2" }
			}`,
		},
		{
			name: "cannot parse previous offset as int",
			out: `{
			  "data": [],
			  "paging": { "next": "?offset=2", "previous": "?offset=foo" }
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux, teardown := setup()
			defer teardown()

			mux.HandleFunc("/anime", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprintf(w, tt.out)
			})

			ctx := context.Background()
			_, _, err := client.Anime.List(ctx, "query")
			if err == nil {
				t.Fatal("Anime.List expected paging error, got no error.")
			}
			if wantPrefix := "paging:"; !strings.HasPrefix(err.Error(), wantPrefix) {
				t.Errorf("Anime.List expected error to start with %q, error is %q", wantPrefix, err.Error())
			}
		})
	}
}

func TestAnimeServiceListError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"mal is down","error":"internal"}`, 500)
	})

	ctx := context.Background()
	_, resp, err := client.Anime.List(ctx, "query")
	if err == nil {
		t.Fatal("Anime.List expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusInternalServerError, "Anime.List")
	testErrorResponse(t, err, ErrorResponse{Message: "mal is down", Err: "internal"})
}

func TestAnimeServiceRanking(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/ranking", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"ranking_type": "all",
			"fields":       "foo,bar",
			"limit":        "10",
			"offset":       "0",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 },
			  "ranking": { "rank": 1 }
		    },
		    {
		      "node": { "id": 2 },
			  "ranking": { "rank": 2 }
		    }
		  ],
		  "paging": {
		    "next": "?offset=4"
		  }
		}`
		fmt.Fprintf(w, out)
	})

	ctx := context.Background()
	got, resp, err := client.Anime.Ranking(ctx, AnimeRankingAll,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("Anime.Ranking returned error: %v", err)
	}
	want := []Anime{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Ranking returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 0, "Anime.Ranking")
}

func TestAnimeServiceSeasonal(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/season/2020/summer", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"sort":   "anime_num_list_users",
			"fields": "foo,bar",
			"limit":  "10",
			"offset": "0",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 }
		    },
		    {
		      "node": { "id": 2 }
		    }
		  ],
		  "paging": {
		    "next": "?offset=4"
		  },
		  "season": {
			"year": 2020,
			"season": "summer"
		  }
		}`
		fmt.Fprintf(w, out)
	})

	ctx := context.Background()
	got, resp, err := client.Anime.Seasonal(ctx, 2020, AnimeSeasonSummer,
		SortSeasonalByAnimeNumListUsers,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("Anime.Seasonal returned error: %v", err)
	}
	want := []Anime{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Seasonal returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 0, "Anime.Seasonal")
}

func TestAnimeServiceSuggested(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/anime/suggestions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"fields": "foo,bar",
			"limit":  "10",
			"offset": "0",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 }
		    },
		    {
		      "node": { "id": 2 }
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
	got, resp, err := client.Anime.Suggested(ctx,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("Anime.Suggested returned error: %v", err)
	}
	want := []Anime{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Suggested returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 2, "Anime.Suggested")
}
