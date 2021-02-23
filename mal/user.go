package mal

import (
	"context"
	"fmt"
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

// MyInfo returns information about the authenticated user.
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

// UserAnimeListOption are options specific to the UserService.AnimeList method.
type UserAnimeListOption func(q *url.Values)

// UserAnimeListStatus filters returned anime list by the status provided. To
// return all anime, don't use this option.
func UserAnimeListStatus(status AnimeStatus) UserAnimeListOption {
	return func(q *url.Values) {
		q.Set("status", string(status))
	}
}

// AnimeStatus is the status of an anime list item.
type AnimeStatus string

// Possible statuses of an nime list item.
const (
	AnimeStatusWatching    AnimeStatus = "watching"
	AnimeStatusCompleted   AnimeStatus = "completed"
	AnimeStatusOnHold      AnimeStatus = "on_hold"
	AnimeStatusDropped     AnimeStatus = "dropped"
	AnimeStatusPlanToWatch AnimeStatus = "plan_to_watch"
)

// FilterAnimeStatus allows to filter the returned anime list by the specified
// status when using the UserService.AnimeList method.
func FilterAnimeStatus(status AnimeStatus) UserAnimeListOption {
	return func(q *url.Values) {
		q.Set("status", string(status))
	}
}

// SortUserAnimeListBy shows the ways the anime results can be sorted.
type SortUserAnimeListBy string

// Possible SortUserAnimeListBy values.
const (
	ByAnimeListScore     SortUserAnimeListBy = "list_score"       // Descending
	ByAnimeListUpdatedAt SortUserAnimeListBy = "list_updated_at"  // Descending
	ByAnimeTitle         SortUserAnimeListBy = "anime_title"      // Ascending
	ByAnimeStartDate     SortUserAnimeListBy = "anime_start_date" // Descending
	ByAnimeID            SortUserAnimeListBy = "anime_id"         // (Under Development) Ascending
)

// SortUserAnimeList allows to choose how the results will be sorted when using
// the UserService.AnimeList method.
func SortUserAnimeList(sort SortUserAnimeListBy) UserAnimeListOption {
	return func(q *url.Values) {
		q.Set("sort", string(sort))
	}
}

// AnimeList gets the anime list of the user indicated by username (or use @me).
// The anime can be sorted and filtered using the SortUserAnimeList and
// FilterAnimeStatus options functions respectively.
func (s *UserService) AnimeList(ctx context.Context, username string, options ...UserAnimeListOption) ([]Anime, *Response, error) {
	oo := make([]func(q *url.Values), len(options))
	for i := range options {
		oo[i] = options[i]
	}
	return s.client.animeList(ctx, fmt.Sprintf("users/%s/animelist", username), oo...)
}
