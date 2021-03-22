package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MangaListOption are options specific to the UserService.MangaList method.
type MangaListOption interface {
	mangaListApply(v *url.Values)
}

// UpdateMyMangaListStatusOption are options specific to the
// MangaService.UpdateMyListStatus method.
type UpdateMyMangaListStatusOption interface {
	updateMyMangaListStatusApply(v *url.Values)
}

func rawOptionFromUpdateMyMangaListStatusOption(o UpdateMyMangaListStatusOption) func(v *url.Values) {
	return func(v *url.Values) {
		o.updateMyMangaListStatusApply(v)
	}
}

// UserManga contains a manga record along with its status on the user's list.
type UserManga struct {
	Manga  Manga           `json:"node"`
	Status MangaListStatus `json:"list_status"`
}

// MangaListStatus shows the status of each manga in a user's manga list.
type MangaListStatus struct {
	Status          MangaStatus `json:"status"`
	IsRereading     bool        `json:"is_rereading"`
	NumVolumesRead  int         `json:"num_volumes_read"`
	NumChaptersRead int         `json:"num_chapters_read"`
	Score           int         `json:"score"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Priority        int         `json:"priority"`
	NumTimesReread  int         `json:"num_times_reread"`
	RereadValue     int         `json:"reread_value"`
	Tags            []string    `json:"tags"`
	Comments        string      `json:"comments"`
}

// mangaList represents the anime list of a user.
type mangaList struct {
	Data   []UserManga `json:"data"`
	Paging Paging      `json:"paging"`
}

func (m mangaList) pagination() Paging { return m.Paging }

// MangaStatus is an option that allows to filter the returned manga list by the
// specified status when using the UserService.MangaList method. It can also be
// passed as an option when updating the manga list.
type MangaStatus string

const (
	// MangaStatusReading returns the manga with status 'reading' from a user's
	// list or sets the status of a list item as such.
	MangaStatusReading MangaStatus = "reading"
	// MangaStatusCompleted returns the manga with status 'completed' from a
	// user's list or sets the status of a list item as such.
	MangaStatusCompleted MangaStatus = "completed"
	// MangaStatusOnHold returns the manga with status 'on hold' from a user's
	// list or sets the status of a list item as such.
	MangaStatusOnHold MangaStatus = "on_hold"
	// MangaStatusDropped returns the manga with status 'dropped' from a user's
	// list or sets the status of a list item as such.
	MangaStatusDropped MangaStatus = "dropped"
	// MangaStatusPlanToRead returns the manga with status 'plan to read' from a
	// user's list or sets the status of a list item as such.
	MangaStatusPlanToRead MangaStatus = "plan_to_read"
)

func (s MangaStatus) mangaListApply(v *url.Values)               { v.Set("status", string(s)) }
func (s MangaStatus) updateMyMangaListStatusApply(v *url.Values) { v.Set("status", string(s)) }

// SortMangaList is an option that sorts the results when getting the user's
// manga list.
type SortMangaList string

const (
	// SortMangaListByListScore sorts results by the score of each item in the
	// list in descending order.
	SortMangaListByListScore SortMangaList = "list_score"
	// SortMangaListByListUpdatedAt sorts results by the most updated entries in
	// the list in descending order.
	SortMangaListByListUpdatedAt SortMangaList = "list_updated_at"
	// SortMangaListByMangaTitle sorts results by the manga title in ascending
	// order.
	SortMangaListByMangaTitle SortMangaList = "manga_title"
	// SortMangaListByMangaStartDate sorts results by the manga start date in
	// descending order.
	SortMangaListByMangaStartDate SortMangaList = "manga_start_date"
	// SortMangaListByMangaID sorts results by the manga ID in ascending order.
	// Note: Currently under development.
	SortMangaListByMangaID SortMangaList = "manga_id"
)

func (s SortMangaList) mangaListApply(v *url.Values) { v.Set("sort", string(s)) }

// MangaList gets the manga list of the user indicated by username (or use @me).
// The manga can be sorted and filtered using the MangaStatus and SortMangaList
// option functions respectively.
func (s *UserService) MangaList(ctx context.Context, username string, options ...MangaListOption) ([]UserManga, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromMangaListOption(options[i])
	}
	list := new(mangaList)
	resp, err := s.client.list(ctx, fmt.Sprintf("users/%s/mangalist", username), list, oo...)
	if err != nil {
		return nil, resp, err
	}
	return list.Data, resp, nil
}

func optionFromMangaListOption(o MangaListOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.mangaListApply(v)
	})
}

// IsRereading is an option that can update if a user is rereading a manga in
// their list.
type IsRereading bool

func (r IsRereading) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("is_rereading", strconv.FormatBool(bool(r)))
}

// NumVolumesRead is an option that can update the number of volumes read of a
// manga in the user's list.
type NumVolumesRead int

func (n NumVolumesRead) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_volumes_read", itoa(int(n)))
}

// NumChaptersRead is an option that can update the number of chapters read of a
// manga in the user's list.
type NumChaptersRead int

func (n NumChaptersRead) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_chapters_read", itoa(int(n)))
}

// NumTimesReread is an option that can update the number of times the user
// has reread a manga in their list.
type NumTimesReread int

func (n NumTimesReread) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("num_times_reread", itoa(int(n)))
}

// RereadValue is an option that can update the reread value of a manga in the
// user's list with values:
//
//     0 = No value
//     1 = Very Low
//     2 = Low
//     3 = Medium
//     4 = High
//     5 = Very High
type RereadValue int

func (r RereadValue) updateMyMangaListStatusApply(v *url.Values) {
	v.Set("reread_value", itoa(int(r)))
}

// UpdateMyListStatus adds the manga specified by mangaID to the user's manga
// list with one or more options added to update the status. If the manga
// already exists in the list, only the status is updated.
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

// DeleteMyListItem deletes a manga from the user's list. If the manga does not
// exist in the user's list, 404 Not Found error is returned.
func (s *MangaService) DeleteMyListItem(ctx context.Context, mangaID int) (*Response, error) {
	u := fmt.Sprintf("manga/%d/my_list_status", mangaID)
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
