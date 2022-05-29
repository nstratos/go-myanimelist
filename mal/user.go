package mal

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// UserService handles communication with the user related methods of the
// MyAnimeList API:
//
// https://myanimelist.net/apiconfig/references/api/v2#tag/user
// https://myanimelist.net/apiconfig/references/api/v2#operation/users_user_id_animelist_get
// https://myanimelist.net/apiconfig/references/api/v2#operation/users_user_id_mangalist_get
type UserService struct {
	client *Client
}

// User represents a MyAnimeList user.
type User struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Picture         string          `json:"picture"`
	Gender          string          `json:"gender"`
	Birthday        string          `json:"birthday"`
	Location        string          `json:"location"`
	JoinedAt        time.Time       `json:"joined_at"`
	AnimeStatistics AnimeStatistics `json:"anime_statistics"`
	TimeZone        string          `json:"time_zone"`
	IsSupporter     bool            `json:"is_supporter"`
}

// AnimeStatistics about the user.
type AnimeStatistics struct {
	NumItemsWatching    int     `json:"num_items_watching"`
	NumItemsCompleted   int     `json:"num_items_completed"`
	NumItemsOnHold      int     `json:"num_items_on_hold"`
	NumItemsDropped     int     `json:"num_items_dropped"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch"`
	NumItems            int     `json:"num_items"`
	NumDaysWatched      float64 `json:"num_days_watched"`
	NumDaysWatching     float64 `json:"num_days_watching"`
	NumDaysCompleted    float64 `json:"num_days_completed"`
	NumDaysOnHold       float64 `json:"num_days_on_hold"`
	NumDaysDropped      float64 `json:"num_days_dropped"`
	NumDays             float64 `json:"num_days"`
	NumEpisodes         int     `json:"num_episodes"`
	NumTimesRewatched   int     `json:"num_times_rewatched"`
	MeanScore           float64 `json:"mean_score"`
}

// MyInfoOption are options specific to the User.MyInfo method.
type MyInfoOption interface {
	myInfoApply(v *url.Values)
}

// MyInfo returns information about the authenticated user.
func (s *UserService) MyInfo(ctx context.Context, options ...MyInfoOption) (*User, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "users/@me")
	if err != nil {
		return nil, nil, err
	}
	q := req.URL.Query()
	for _, o := range options {
		o.myInfoApply(&q)
	}
	req.URL.RawQuery = q.Encode()

	u := new(User)
	resp, err := s.client.Do(ctx, req, u)
	if err != nil {
		return nil, resp, err
	}

	return u, resp, nil
}
