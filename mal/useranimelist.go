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

// AnimeListOption are options specific to the UserService.AnimeList method.
type AnimeListOption interface {
	animeListApply(v *url.Values)
}

// AnimeStatus is an option that allows to filter the returned anime list by the
// specified status when using the UserService.AnimeList method. It can also be
// passed as an option when updating the anime list.
type AnimeStatus string

const (
	// AnimeStatusWatching returns the anime with status 'watching' from a
	// user's list or sets the status of a list item as such.
	AnimeStatusWatching AnimeStatus = "watching"
	// AnimeStatusCompleted returns the anime with status 'completed' from a
	// user's list or sets the status of a list item as such.
	AnimeStatusCompleted AnimeStatus = "completed"
	// AnimeStatusOnHold returns the anime with status 'on hold' from a user's
	// list or sets the status of a list item as such.
	AnimeStatusOnHold AnimeStatus = "on_hold"
	// AnimeStatusDropped returns the anime with status 'dropped' from a user's
	// list or sets the status of a list item as such.
	AnimeStatusDropped AnimeStatus = "dropped"
	// AnimeStatusPlanToWatch returns the anime with status 'plan to watch' from
	// a user's list or sets the status of a list item as such.
	AnimeStatusPlanToWatch AnimeStatus = "plan_to_watch"
)

func (s AnimeStatus) animeListApply(v *url.Values)               { v.Set("status", string(s)) }
func (s AnimeStatus) updateMyAnimeListStatusApply(v *url.Values) { v.Set("status", string(s)) }

// SortAnimeList is an option that sorts the results when getting the user's
// anime list.
type SortAnimeList string

const (
	// SortAnimeListByListScore sorts results by the score of each item in the
	// list in descending order.
	SortAnimeListByListScore SortAnimeList = "list_score"
	// SortAnimeListByListUpdatedAt sorts results by the most updated entries in
	// the list in descending order.
	SortAnimeListByListUpdatedAt SortAnimeList = "list_updated_at"
	// SortAnimeListByAnimeTitle sorts results by the anime title in ascending
	// order.
	SortAnimeListByAnimeTitle SortAnimeList = "anime_title"
	// SortAnimeListByAnimeStartDate sorts results by the anime start date in
	// descending order.
	SortAnimeListByAnimeStartDate SortAnimeList = "anime_start_date"
	// SortAnimeListByAnimeID sorts results by the anime ID in ascending order.
	// Note: Currently under development.
	SortAnimeListByAnimeID SortAnimeList = "anime_id"
)

func (s SortAnimeList) animeListApply(v *url.Values) { v.Set("sort", string(s)) }

// UserAnime contains an anime record along with its status on the user's list.
type UserAnime struct {
	Anime  Anime           `json:"node"`
	Status AnimeListStatus `json:"list_status"`
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

// animeList represents the anime list of a user.
type animeList struct {
	Data   []UserAnime `json:"data"`
	Paging Paging      `json:"paging"`
}

func (a animeList) pagination() Paging { return a.Paging }

// AnimeList gets the anime list of the user indicated by username (or use @me).
// The anime can be sorted and filtered using the AnimeStatus and SortAnimeList
// option functions respectively.
func (s *UserService) AnimeList(ctx context.Context, username string, options ...AnimeListOption) ([]UserAnime, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromAnimeListOption(options[i])
	}
	list := new(animeList)
	resp, err := s.client.list(ctx, fmt.Sprintf("users/%s/animelist", username), list, oo...)
	if err != nil {
		return nil, resp, err
	}
	return list.Data, resp, nil
}

func optionFromAnimeListOption(o AnimeListOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.animeListApply(v)
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

// NumEpisodesWatched is an option that can update the number of episodes
// watched of an anime in the user's list.
type NumEpisodesWatched int

func (n NumEpisodesWatched) updateMyAnimeListStatusApply(v *url.Values) {
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
// the user's list with values:
//
//     0 = No value
//     1 = Very Low
//     2 = Low
//     3 = Medium
//     4 = High
//     5 = Very High
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

// DeleteMyListItem deletes an anime from the user's list. If the anime does not
// exist in the user's list, 404 Not Found error is returned.
func (s *AnimeService) DeleteMyListItem(ctx context.Context, animeID int) (*Response, error) {
	u := fmt.Sprintf("anime/%d/my_list_status", animeID)
	req, err := s.client.NewRequest(http.MethodDelete, u)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
