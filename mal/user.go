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
	req, err := s.client.NewRequest(http.MethodGet, "users/@me")
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

// AnimeListStatus shows the status of each anime in a user anime list.
type AnimeListStatus struct {
	Status             AnimeStatus `json:"status"`
	Score              int         `json:"score"`
	NumEpisodesWatched int         `json:"num_episodes_watched"`
	IsRewatching       bool        `json:"is_rewatching"`
	UpdatedAt          *time.Time  `json:"updated_at"`
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

type Score int

func (s Score) updateMyAnimeListStatusApply(v *url.Values) { v.Set("score", itoa(int(s))) }
func (s Score) updateMyMangaListStatusApply(v *url.Values) { v.Set("score", itoa(int(s))) }

type NumWatchedEpisodes int

func (n NumWatchedEpisodes) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("num_watched_episodes", itoa(int(n)))
}

type NumTimesRewatched int

func (n NumTimesRewatched) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("num_times_rewatched", itoa(int(n)))
}

type IsRewatching bool

func (r IsRewatching) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("is_rewatching", strconv.FormatBool(bool(r)))
}

type RewatchValue int

func (r RewatchValue) updateMyAnimeListStatusApply(v *url.Values) {
	v.Set("rewatch_value", itoa(int(r)))
}

type Priority int

func (p Priority) updateMyAnimeListStatusApply(v *url.Values) { v.Set("priority", itoa(int(p))) }
func (p Priority) updateMyMangaListStatusApply(v *url.Values) { v.Set("priority", itoa(int(p))) }

type Tags []string

func (t Tags) updateMyAnimeListStatusApply(v *url.Values) { v.Set("tags", strings.Join(t, ",")) }
func (t Tags) updateMyMangaListStatusApply(v *url.Values) { v.Set("tags", strings.Join(t, ",")) }

type Comments string

func (c Comments) updateMyAnimeListStatusApply(v *url.Values) { v.Set("comments", string(c)) }
func (c Comments) updateMyMangaListStatusApply(v *url.Values) { v.Set("comments", string(c)) }

// UpdateMyListStatus adds anime specified by the animeID to the users anime
// list with the status specified by animeStatus. If the anime already exists in
// the list, only the status is updated.
//
// This endpoint updates only values specified by the parameter.
//
// TODO(nstratos): How does it work with Go's default values?
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

// UpdateMyMangaListStatusOption are options specific to the
// AmangaService.UpdateMyListStatus method.
type UpdateMyMangaListStatusOption interface {
	updateMyMangaListStatusApply(v *url.Values)
}

func rawOptionFromUpdateMyMangaListStatusOption(o UpdateMyMangaListStatusOption) func(v *url.Values) {
	return func(v *url.Values) {
		o.updateMyMangaListStatusApply(v)
	}
}

type MangaListStatus struct {
	Status          string        `json:"status"`
	IsRereading     bool          `json:"is_rereading"`
	NumVolumesRead  int           `json:"num_volumes_read"`
	NumChaptersRead int           `json:"num_chapters_read"`
	Score           int           `json:"score"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Priority        int           `json:"priority"`
	NumTimesReread  int           `json:"num_times_reread"`
	RereadValue     int           `json:"reread_value"`
	Tags            []interface{} `json:"tags"`
	Comments        string        `json:"comments"`
}

// MangaStatus is an option that allows to filter the returned anime list by the
// specified status when using the UserService.MangaList method. It can also be
// passed as an option when updating the manga list.
type MangaStatus string

// Possible statuses of an nime list item.
const (
	MangaStatusReading    MangaStatus = "reading"
	MangaStatusCompleted  MangaStatus = "completed"
	MangaStatusOnHold     MangaStatus = "on_hold"
	MangaStatusDropped    MangaStatus = "dropped"
	MangaStatusPlanToRead MangaStatus = "plan_to_read"
)

func (s MangaStatus) mangaListApply(v *url.Values)               { v.Set("status", string(s)) }
func (s MangaStatus) updateMyMangaListStatusApply(v *url.Values) { v.Set("status", string(s)) }

type IsRereading bool

func (r IsRereading) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("is_rereading", strconv.FormatBool(bool(r)))
}

type NumVolumesRead int

func (n NumVolumesRead) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_volumes_read", itoa(int(n)))
}

type NumChaptersRead int

func (n NumChaptersRead) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_chapters_read", itoa(int(n)))
}

type NumTimesReread int

func (n NumTimesReread) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_times_reread", itoa(int(n)))
}

type RereadValue int

func (r RereadValue) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("reread_value", itoa(int(r)))
}

// UpdateMyListStatus adds manga specified by the mangaID to the users anime
// list with the status specified by animeStatus. If the anime already exists in
// the list, only the status is updated.
//
// This endpoint updates only values specified by the parameter.
func (s *MangaService) UpdateMyListStatus(ctx context.Context, mangaID int, options ...UpdateMyMangaListStatusOption) (*MangaListStatus, *Response, error) {
	u := fmt.Sprintf("manga/%d/my_list_status", mangaID)
	rawOptions := make([]func(v *url.Values), len(options))
	for i := range options {
		rawOptions[i] = rawOptionFromUpdateMyMangaListStatusOption(options[i])
	}
	req, err := s.client.NewRequest(http.MethodPatch, u, rawOptions...)
	if err != nil {
		return nil, nil, err
	}

	m := new(MangaListStatus)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}
