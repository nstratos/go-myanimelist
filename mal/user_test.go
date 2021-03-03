package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestUserMarshal(t *testing.T) {
	testJSONMarshal(t, &User{}, "{}")

	u := &User{
		ID:       6548478,
		Name:     "rin-0911-3",
		Gender:   "g",
		Location: "l",
		Picture:  "p",
		JoinedAt: time.Date(2017, time.September, 11, 10, 27, 46, 0, time.UTC),
		AnimeStatistics: AnimeStatistics{
			NumItemsWatching:    1,
			NumItemsCompleted:   2,
			NumItemsOnHold:      3,
			NumItemsDropped:     4,
			NumItemsPlanToWatch: 5,
			NumItems:            6,
			NumDaysWatched:      7.1,
			NumDaysWatching:     8.2,
			NumDaysCompleted:    9.3,
			NumDaysOnHold:       10.4,
			NumDaysDropped:      11.5,
			NumDays:             12.6,
			NumEpisodes:         13,
			NumTimesRewatched:   14,
			MeanScore:           15.7,
		},
	}
	want := `{
		"id": 6548478,
		"name": "rin-0911-3",
		"gender": "g",
		"location": "l",
		"picture": "p",
		"joined_at": "2017-09-11T10:27:46+00:00",
		"anime_statistics": {
			"num_items_watching": 1,
			"num_items_completed": 2,
			"num_items_on_hold": 3,
			"num_items_dropped": 4,
			"num_items_plan_to_watch": 5,
			"num_items": 6,
			"num_days_watched": 7.1,
			"num_days_watching": 8.2,
			"num_days_completed": 9.3,
			"num_days_on_hold": 10.4,
			"num_days_dropped": 11.5,
			"num_days": 12.6,
			"num_episodes": 13,
			"num_times_rewatched": 14,
			"mean_score": 15.7
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestUserServiceMyInfo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/@me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"id":1}`)
	})

	ctx := context.Background()
	u, _, err := client.User.MyInfo(ctx)
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
		testMethod(t, r, "GET")
		http.Error(w, `{"message":"","error":"not_found"}`, 404)
	})

	ctx := context.Background()
	_, _, err := client.User.MyInfo(ctx)
	if err == nil {
		t.Fatal("User.MyInfo expected not found error, got no error.")
	}
	testErrorResponse(t, err, ErrorResponse{Err: "not_found"})
}
