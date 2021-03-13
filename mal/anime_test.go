package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestAnimeMarshal(t *testing.T) {
	testJSONMarshal(t, &Anime{}, "{}")
	createdAt := time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC)
	u := &Anime{
		ID:    1,
		Title: "t",
		MainPicture: Picture{
			Medium: "m",
			Large:  "l",
		},
		AlternativeTitles: Titles{
			Synonyms: []string{"s1", "s2"},
			En:       "e",
			Ja:       "j",
		},
		StartDate:       "2020-01-01",
		EndDate:         "2021-01-01",
		Synopsis:        "s",
		Mean:            1.1,
		Rank:            1,
		Popularity:      1,
		NumListUsers:    1,
		NumScoringUsers: 1,
		NSFW:            "white",
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		MediaType:       "m",
		Status:          "s",
		Genres:          []Genre{{ID: 1, Name: "n"}},
		MyListStatus: AnimeListStatus{
			Status:             "s",
			Score:              1,
			NumEpisodesWatched: 1,
			IsRewatching:       true,
			UpdatedAt:          updatedAt,
			Priority:           1,
			NumTimesRewatched:  1,
			RewatchValue:       1,
			Tags:               []string{"t1", "t2"},
			Comments:           "c",
		},
		NumEpisodes:            1,
		StartSeason:            StartSeason{Year: 1, Season: "s"},
		Broadcast:              Broadcast{DayOfTheWeek: "d", StartTime: "15:00"},
		Source:                 "s",
		AverageEpisodeDuration: 1,
		Rating:                 "r",
		Pictures:               []Picture{},
		Background:             "b",
		RelatedAnime:           []RelatedAnime{{Node: Anime{ID: 1}, RelationType: "r", RelationTypeFormatted: "r"}},
		RelatedManga:           []RelatedManga{},
		Recommendations:        []RecommendedAnime{},
		Studios:                []Studio{{ID: 1, Name: "n"}, {ID: 2, Name: "n"}},
		Statistics: Statistics{
			Status: Status{
				Watching:    "1",
				Completed:   "1",
				OnHold:      "1",
				Dropped:     "1",
				PlanToWatch: "1",
			},
			NumListUsers: 1,
		},
	}
	want := `
	{
		"id": 1,
		"title": "t",
		"main_picture": {
		  "medium": "m",
		  "large": "j"
		},
		"alternative_titles": {
		  "synonyms": ["s1", "s2"],
		  "en": "e",
		  "ja": "j"
		},
		"start_date": "2020-01-01",
		"end_date": "2021-01-01",
		"synopsis": "s",
		"mean": 1.1,
		"rank": 1,
		"popularity": 1,
		"num_list_users": 1,
		"num_scoring_users": 1,
		"nsfw": "white",
		"created_at": "2020-01-01T15:00:00Z",
		"updated_at": "2020-01-01T15:00:00Z",
		"media_type": "m",
		"status": "s",
		"genres": [{"id": 1, "name": "n"}],
		"my_list_status": {
		  "status": "s",
		  "score": 1,
		  "priority": 1,
		  "num_episodes_watched": 1,
		  "num_times_rewatched": 1,
		  "is_rewatching": true,
		  "updated_at": "2020-01-01T15:00:00Z",
		  "tags": ["t1", "t2"],
		  "comments": "c"
		},
		"num_episodes": 1,
		"start_season": {
		  "year": 2020,
		  "season": "s"
		},
		"broadcast": {
		  "day_of_the_week": "d",
		  "start_time": "15:00"
		},
		"source": "s",
		"average_episode_duration": 1,
		"rating": "r",
		"pictures": [
		  {
			"medium": "m",
			"large": "l"
		  }
		],
		"background": "b",
		"related_anime": [
		  {
			"node": {
			  "id": 1,
			  "title": "t",
			  "main_picture": {
				"medium": "m",
				"large": "l"
			  }
			},
			"relation_type": "r",
			"relation_type_formatted": "r"
		  }
		],
		"related_manga": [],
		"recommendations": [
		  {
			"node": {
			  "id": 1,
			  "title": "t",
			  "main_picture": {
				"medium": "m",
				"large": "l"
			  }
			},
			"num_recommendations": 4
		  }
		],
		"studios": [{ "id": 1, "name": "n" }],
		"statistics": {
		  "status": {
			"watching": "1",
			"completed": "1",
			"on_hold": "1",
			"dropped": "1",
			"plan_to_watch": "1"
		  },
		  "num_list_users": 1
		}
	  }`
	testJSONMarshal(t, u, want)
}

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
