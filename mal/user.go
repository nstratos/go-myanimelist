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

// animeListOption are options specific to the UserService.AnimeList method.
type animeListOption interface {
	animeListApply(v *url.Values)
}

// AnimeStatus is an option that allows to filter the returned anime list by the
// specified status when using the UserService.AnimeList method.
type AnimeStatus string

// Possible statuses of an nime list item.
const (
	AnimeStatusWatching    AnimeStatus = "watching"
	AnimeStatusCompleted   AnimeStatus = "completed"
	AnimeStatusOnHold      AnimeStatus = "on_hold"
	AnimeStatusDropped     AnimeStatus = "dropped"
	AnimeStatusPlanToWatch AnimeStatus = "plan_to_watch"
)

func (s AnimeStatus) animeListApply(v *url.Values) { v.Set("status", string(s)) }

// SortAnimeList is an option that sorts the results returned by the
// UserService.AnimeList method.
type SortAnimeList string

// Possible sorting values.
const (
	SortAnimeListByAnimeListScore     SortAnimeList = "list_score"       // Descending
	SortAnimeListByAnimeListUpdatedAt SortAnimeList = "list_updated_at"  // Descending
	SortAnimeListByAnimeTitle         SortAnimeList = "anime_title"      // Ascending
	SortAnimeListByAnimeStartDate     SortAnimeList = "anime_start_date" // Descending
	SortAnimeListByAnimeID            SortAnimeList = "anime_id"         // (Under Development) Ascending
)

func (s SortAnimeList) animeListApply(v *url.Values) { v.Set("sort", string(s)) }

// AnimeWithStatus contains an anime record along with its list status.
type AnimeWithStatus struct {
	Anime  Anime
	Status AnimeListStatus
}

// AnimeList gets the anime list of the user indicated by username (or use @me).
// The anime can be sorted and filtered using the AnimeStatus and SortAnimeList
// options functions respectively.
//
// Example:
//
//     anime, _, err := c.User.AnimeList(ctx, "leonteus",
//         mal.Limit(10),
//         mal.Fields{"rank", "popularity"},
//         mal.SortAnimeListByAnimeListScore,
//     )
//     if err != nil {
//         return err
//     }
//     for _, a := range anime {
//         fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
//     }
func (s *UserService) AnimeList(ctx context.Context, username string, options ...animeListOption) ([]AnimeWithStatus, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromAnimeListOption(options[i])
	}
	return s.animeListWithStatus(ctx, fmt.Sprintf("users/%s/animelist", username), oo...)
}

func optionFromAnimeListOption(o animeListOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.animeListApply(v)
	})
}

func (s *UserService) animeListWithStatus(ctx context.Context, path string, options ...Option) ([]AnimeWithStatus, *Response, error) {
	list, resp, err := s.client.animeList(ctx, path, options...)
	if err != nil {
		return nil, resp, err
	}
	anime := make([]AnimeWithStatus, len(list.Data))
	for i := range list.Data {
		anime[i].Anime = list.Data[i].Anime
		anime[i].Status = list.Data[i].Status
	}
	return anime, resp, nil
}
