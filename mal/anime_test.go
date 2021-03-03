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
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	a, _, err := client.Anime.Details(ctx, 1)
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
	_, resp, err := client.Anime.List(ctx, "query",
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
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
	got, resp, err := client.Anime.Ranking(ctx, RankingAiring,
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

// func TestAnimeService_Add(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/animelist/add/", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")
// 		testID(t, r, "55")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testContentType(t, r, "application/x-www-form-urlencoded")
// 		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, 1))
// 		fmt.Fprintf(w, "Created")
// 	})

// 	_, err := client.Anime.Add(55, AnimeEntry{Status: Current})
// 	if err != nil {
// 		t.Errorf("Anime.Add returned error %v", err)
// 	}
// }

// func TestAnimeService_Update(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/animelist/update/", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")
// 		testID(t, r, "55")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testContentType(t, r, "application/x-www-form-urlencoded")
// 		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, 3))
// 		fmt.Fprintf(w, "Updated")
// 	})

// 	_, err := client.Anime.Update(55, AnimeEntry{Status: OnHold})
// 	if err != nil {
// 		t.Errorf("Anime.Update returned error %v", err)
// 	}
// }

// func TestAnimeService_Update_invalidID(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/animelist/update/", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")
// 		testID(t, r, "0")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testContentType(t, r, "application/x-www-form-urlencoded")
// 		testFormValue(t, r, "data", fmt.Sprintf(zeroEntry, 3))
// 		http.Error(w, "Invalid ID", http.StatusNotImplemented)
// 	})

// 	response, err := client.Anime.Update(0, AnimeEntry{Status: OnHold})

// 	if err == nil {
// 		t.Errorf("Anime.Update invalid ID should return err")
// 	}

// 	if response == nil {
// 		t.Errorf("Anime.Update invalid ID should return also return response")
// 	}
// }

// func TestAnimeService_Search(t *testing.T) {
// 	client, mux, teardown := setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/anime/search.xml", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testURLValues(t, r, urlValues{"q": "query"})
// 		fmt.Fprintf(w, `
// 			<anime>
// 				<entry>
// 					<title>title1</title>
// 					<id>55</id>
// 				</entry>
// 				<entry>
// 					<title>title2</title>
// 					<id>56</id>
// 				</entry>
// 			</anime>`)
// 	})

// 	result, _, err := client.Anime.Search("query")
// 	if err != nil {
// 		t.Errorf("Anime.Search returned error %v", err)
// 	}
// 	want := &AnimeResult{
// 		[]AnimeRow{
// 			{ID: 55, Title: "title1"},
// 			{ID: 56, Title: "title2"},
// 		},
// 	}
// 	if !reflect.DeepEqual(result, want) {
// 		t.Errorf("Anime.Search returned %+v, want %+v", result, want)
// 	}
// }

// func TestAnimeService_Search_noContent(t *testing.T) {
// 	client, mux, teardown := setup()
// 	defer teardown()

// 	mux.HandleFunc("/api/anime/search.xml", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testBasicAuth(t, r, true, "TestUser", "TestPass")
// 		testURLValues(t, r, urlValues{"q": "foo"})
// 		http.Error(w, "no content", http.StatusNoContent)
// 	})

// 	result, _, err := client.Anime.Search("foo")

// 	if err == nil {
// 		t.Errorf("Anime.Search for non existent query expected to return err")
// 	}

// 	if got, want := err, ErrNoContent; got != want {
// 		t.Errorf("Anime.Search for non existent query returned err %v, want %v", got, want)
// 	}

// 	if got := result; got != nil {
// 		t.Errorf("Anime.Search for non existent query returned result = %v, want %v", got, nil)
// 	}
// }

// func TestAnimeService_List(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	mux.HandleFunc("/malappinfo.php", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testBasicAuth(t, r, false, "", "")
// 		testURLValues(t, r, urlValues{
// 			"status": "all",
// 			"type":   "anime",
// 			"u":      "AnotherTestUser",
// 		})
// 		fmt.Fprintf(w, `
// 			<myanimelist>
// 				<myinfo>
// 					<user_id>56</user_id>
// 					<user_name>AnotherTestUser</user_name>
// 				</myinfo>
// 				<anime>
// 					<series_animedb_id>1</series_animedb_id>
// 					<series_title>series title</series_title>
// 					<my_id>1234</my_id>
// 					<my_status>3</my_status>
// 				</anime>
// 			</myanimelist>
// 			`)
// 	})

// 	result, _, err := client.Anime.List("AnotherTestUser")
// 	if err != nil {
// 		t.Errorf("Anime.List returned error %v", err)
// 	}
// 	want := &AnimeList{
// 		MyInfo: AnimeMyInfo{ID: 56, Name: "AnotherTestUser"},
// 		Anime: []Anime2{
// 			{
// 				SeriesAnimeDBID: 1,
// 				SeriesTitle:     "series title",
// 				MyID:            1234,
// 				MyStatus:        3,
// 			},
// 		},
// 	}
// 	if !reflect.DeepEqual(result, want) {
// 		t.Errorf("Anime.List returned %+v, want %+v", result, want)
// 	}
// }

// func TestAnimeService_List_invalidUsername(t *testing.T) {
// 	client, mux, teardown := setup()
// 	defer teardown()

// 	mux.HandleFunc("/malappinfo.php", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testBasicAuth(t, r, false, "", "")
// 		testURLValues(t, r, urlValues{
// 			"status": "all",
// 			"type":   "anime",
// 			"u":      "InvalidUser",
// 		})
// 		fmt.Fprintf(w, `
// 			<myanimelist>
// 				<error>Invalid username</error>
// 			</myanimelist>
// 			`)
// 	})

// 	ctx := context.Background()
// 	result, _, err := client.Anime.List(ctx, "InvalidUser", 100, 0)

// 	if err == nil {
// 		t.Errorf("Anime.List for invalid user expected to return err")
// 	}

// 	want := &animeList{Error: "Invalid username"}
// 	if !reflect.DeepEqual(want, result) {
// 		t.Errorf("Anime.List for invalid user returned result = %v, want %v", result, want)
// 	}
// }

// func TestAnimeService_List_httpError(t *testing.T) {
// 	client, mux, teardown := setup()
// 	defer teardown()

// 	mux.HandleFunc("/malappinfo.php", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testBasicAuth(t, r, false, "", "")
// 		testURLValues(t, r, urlValues{
// 			"status": "all",
// 			"type":   "anime",
// 			"u":      "TestUser",
// 		})
// 		http.Error(w, "something broke", http.StatusInternalServerError)
// 	})

// 	ctx := context.Background()
// 	result, _, err := client.Anime.List(ctx, "TestUser")

// 	if err == nil {
// 		t.Errorf("Anime.List for server error expected to return err")
// 	}

// 	if got := result; got != nil {
// 		t.Errorf("Anime.List for server error returned result = %v, want %v", got, nil)
// 	}
// }
