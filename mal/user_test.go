package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestUserServiceMyInfo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/@me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"fields": "time_zone,is_supporter",
		})
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	u, _, err := client.User.MyInfo(ctx,
		Fields{"time_zone", "is_supporter"},
	)
	if err != nil {
		t.Errorf("User.MyInfo returned error: %v", err)
	}
	want := &User{ID: 1}
	if got := u; !reflect.DeepEqual(got, want) {
		t.Errorf("User.MyInfo returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestUserServiceMyInfoError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/@me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"","error":"not_found"}`, 404)
	})

	ctx := context.Background()
	_, _, err := client.User.MyInfo(ctx)
	if err == nil {
		t.Fatal("User.MyInfo expected not found error, got no error.")
	}
	testErrorResponse(t, err, ErrorResponse{Err: "not_found"})
}

func TestUserMarshal(t *testing.T) {
	testJSONMarshal(t, &User{}, "{}")
	u := &User{
		ID:       1,
		Name:     "foo",
		Picture:  "https://api-cdn.myanimelist.net/images/userimages/1.jpg?t=1653811800",
		Gender:   "baz",
		Birthday: "1990-01-01",
		Location: "bar",
		JoinedAt: time.Date(2015, 8, 20, 11, 11, 55, 0, time.UTC),
		AnimeStatistics: AnimeStatistics{
			NumItemsWatching:    1,
			NumItemsCompleted:   1,
			NumItemsOnHold:      1,
			NumItemsDropped:     1,
			NumItemsPlanToWatch: 1,
			NumItems:            1,
			NumDaysWatched:      1.1,
			NumDaysWatching:     1.1,
			NumDaysCompleted:    1.1,
			NumDaysOnHold:       1.1,
			NumDaysDropped:      1.1,
			NumDays:             1.1,
			NumEpisodes:         1,
			NumTimesRewatched:   1,
			MeanScore:           1.1,
		},
		TimeZone:    "America/Guyana",
		IsSupporter: true,
	}
	want := `{
		"id": 1,
		"name": "foo",
		"picture": "https://api-cdn.myanimelist.net/images/userimages/1.jpg?t=1653811800",
		"gender": "baz",
		"birthday": "1990-01-01",
		"location": "bar",
		"joined_at": "2015-08-20T11:11:55Z",
		"anime_statistics": {
			"num_items_watching": 1,
			"num_items_completed": 1,
			"num_items_on_hold": 1,
			"num_items_dropped": 1,
			"num_items_plan_to_watch": 1,
			"num_items": 1,
			"num_days_watched": 1.1,
			"num_days_watching": 1.1,
			"num_days_completed": 1.1,
			"num_days_on_hold": 1.1,
			"num_days_dropped": 1.1,
			"num_days": 1.1,
			"num_episodes": 1,
			"num_times_rewatched": 1,
			"mean_score": 1.1
		},
		"time_zone": "America/Guyana",
		"is_supporter": true
	}`
	testJSONMarshal(t, u, want)
}
