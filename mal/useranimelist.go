package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// animeListOption are options specific to the UserService.AnimeList method.
type animeListOption interface {
	animeListApply(v *url.Values)
}

// AnimeStatus is an option that allows to filter the returned anime list by the
// specified status when using the UserService.AnimeList method.
type AnimeStatus string

// Possible statuses of an anime in the user's list.
const (
	AnimeStatusWatching    AnimeStatus = "watching"
	AnimeStatusCompleted   AnimeStatus = "completed"
	AnimeStatusOnHold      AnimeStatus = "on_hold"
	AnimeStatusDropped     AnimeStatus = "dropped"
	AnimeStatusPlanToWatch AnimeStatus = "plan_to_watch"
)

func (s AnimeStatus) animeListApply(v *url.Values)               { v.Set("status", string(s)) }
func (s AnimeStatus) updateMyAnimeListStatusApply(v *url.Values) { v.Set("status", string(s)) }

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

// AnimeListStatus shows the status of each anime in a user's anime list.
type AnimeListStatus struct {
	Status             AnimeStatus `json:"status"`
	Score              int         `json:"score"`
	NumEpisodesWatched int         `json:"num_episodes_watched"`
	IsRewatching       bool        `json:"is_rewatching"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Priority           int         `json:"priority"`
	NumTimesRewatched  int         `json:"num_times_rewatched"`
	RewatchValue       int         `json:"rewatch_value"`
	Tags               []string    `json:"tags"`
	Comments           string      `json:"comments"`
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

// UpdateAnimeListStatus shows the status of each anime in a user anime list.
type UpdateAnimeListStatus struct {
	Status             *AnimeStatus
	Score              *int
	NumWatchedEpisodes *int
	IsRewatching       *bool
	UpdatedAt          *time.Time
	Priority           *int
	NumTimesRewatched  *int
	RewatchValue       *int
	Tags               []string
	Comments           *string
}

func optionFromUpdateAnimeListStatus(u UpdateAnimeListStatus) optionFunc {
	return optionFunc(func(v *url.Values) {
		if u.Status != nil {
			v.Set("status", string(*u.Status))
		}
		if u.NumWatchedEpisodes != nil {
			v.Set("num_watched_episodes", strconv.Itoa(*u.NumWatchedEpisodes))
		}
	})
}

// UpdateMyAnimeListStatusOption are options specific to the
// AnimeService.UpdateMyListStatus method.
type UpdateMyAnimeListStatusOption interface {
	updateMyAnimeListStatusApply(v *url.Values)
}

func rawOptionFromUpdateMyAnimeListStatusOption(o UpdateMyAnimeListStatusOption) func(v *url.Values) {
	return func(v *url.Values) {
		o.updateMyAnimeListStatusApply(v)
	}
}

var itoa = strconv.Itoa

// Score is an option that can update the anime and manga list scores with
// values 0-10.
type Score int

func (s Score) updateMyAnimeListStatusApply(v *url.Values) { v.Set("score", itoa(int(s))) }
func (s Score) updateMyMangaListStatusApply(v *url.Values) { v.Set("score", itoa(int(s))) }

// NumWatchedEpisodes is an option that can update the number of episodes
// watched of an anime in the user's list.
type NumWatchedEpisodes int

func (n NumWatchedEpisodes) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("num_watched_episodes", itoa(int(n)))
}

// NumTimesRewatched is an option that can update the number of times the user
// has rewatched an anime in their list.
type NumTimesRewatched int

func (n NumTimesRewatched) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("num_times_rewatched", itoa(int(n)))
}

// IsRewatching is an option that can update if a user is rewatching an anime in
// their list.
type IsRewatching bool

func (r IsRewatching) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("is_rewatching", strconv.FormatBool(bool(r)))
}

// RewatchValue is an option that can update the rewatch value of an anime in
// the user's list with values 0-5.
type RewatchValue int

func (r RewatchValue) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("rewatch_value", itoa(int(r)))
}

// Priority is an option that allows to update the priority of an anime or manga
// in the user's list with values 0=Low, 1=Medium, 2=High.
type Priority int

func (p Priority) updateMyAnimeListStatusApply(v *url.Values) { v.Set("priority", itoa(int(p))) }
func (p Priority) updateMyMangaListStatusApply(v *url.Values) { v.Set("priority", itoa(int(p))) }

// Tags is an option that allows to update the comma-separated tags of anime and
// manga in the user's list.
type Tags []string

func (t Tags) updateMyAnimeListStatusApply(v *url.Values) { v.Set("tags", strings.Join(t, ",")) }
func (t Tags) updateMyMangaListStatusApply(v *url.Values) { v.Set("tags", strings.Join(t, ",")) }

// Comments is an option that allows to update the comment of anime and manga in
// the user's list.
type Comments string

func (c Comments) updateMyAnimeListStatusApply(v *url.Values) { v.Set("comments", string(c)) }
func (c Comments) updateMyMangaListStatusApply(v *url.Values) { v.Set("comments", string(c)) }

// UpdateMyListStatus adds the anime specified by animeID to the user's anime
// list with one or more options added to update the status. If the anime
// already exists in the list, only the status is updated.
func (s *AnimeService) UpdateMyListStatus(ctx context.Context, animeID int, options ...UpdateMyAnimeListStatusOption) (*AnimeListStatus, *Response, error) {
	u := fmt.Sprintf("anime/%d/my_list_status", animeID)
	rawOptions := make([]func(v *url.Values), len(options))
	for i := range options {
		rawOptions[i] = rawOptionFromUpdateMyAnimeListStatusOption(options[i])
	}
	req, err := s.client.NewRequest(http.MethodPatch, u, rawOptions...)
	if err != nil {
		return nil, nil, err
	}

	a := new(AnimeListStatus)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}
