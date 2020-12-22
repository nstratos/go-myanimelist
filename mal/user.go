package mal

import (
	"context"
	"net/http"
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
	ID              int64           `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Gender          string          `json:"gender,omitempty"`
	Location        string          `json:"location,omitempty"`
	Picture         string          `json:"picture,omitempty"`
	JoinedAt        time.Time       `json:"joined_at,omitempty"`
	AnimeStatistics AnimeStatistics `json:"anime_statistics,omitempty"`
}

// AnimeStatistics about the user.
type AnimeStatistics struct {
	NumItemsWatching    int     `json:"num_items_watching,omitempty"`
	NumItemsCompleted   int     `json:"num_items_completed,omitempty"`
	NumItemsOnHold      int     `json:"num_items_on_hold,omitempty"`
	NumItemsDropped     int     `json:"num_items_dropped,omitempty"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch,omitempty"`
	NumItems            int     `json:"num_items,omitempty"`
	NumDaysWatched      float64 `json:"num_days_watched,omitempty"`
	NumDaysWatching     float64 `json:"num_days_watching,omitempty"`
	NumDaysCompleted    float64 `json:"num_days_completed,omitempty"`
	NumDaysOnHold       float64 `json:"num_days_on_hold,omitempty"`
	NumDaysDropped      float64 `json:"num_days_dropped,omitempty"`
	NumDays             float64 `json:"num_days,omitempty"`
	NumEpisodes         int     `json:"num_episodes,omitempty"`
	NumTimesRewatched   int     `json:"num_times_rewatched,omitempty"`
	MeanScore           float64 `json:"mean_score,omitempty"`
}

// MyInfo returns information about the authorized user.
func (s *UserService) MyInfo(ctx context.Context) (*User, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "users/@me", nil)
	if err != nil {
		return nil, nil, err
	}

	u := new(User)
	resp, err := s.client.Do(ctx, req, u)
	if err != nil {
		return nil, resp, err
	}

	return u, resp, nil
}
