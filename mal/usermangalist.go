package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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

// MangaListStatus shows the status of each manga in a user's manga list.
type MangaListStatus struct {
	Status          string    `json:"status"`
	IsRereading     bool      `json:"is_rereading"`
	NumVolumesRead  int       `json:"num_volumes_read"`
	NumChaptersRead int       `json:"num_chapters_read"`
	Score           int       `json:"score"`
	UpdatedAt       time.Time `json:"updated_at"`
	Priority        int       `json:"priority"`
	NumTimesReread  int       `json:"num_times_reread"`
	RereadValue     int       `json:"reread_value"`
	Tags            []string  `json:"tags"`
	Comments        string    `json:"comments"`
}

// MangaStatus is an option that allows to filter the returned manga list by the
// specified status when using the UserService.MangaList method. It can also be
// passed as an option when updating the manga list.
type MangaStatus string

// Possible statuses of a manga in the user's list.
const (
	MangaStatusReading    MangaStatus = "reading"
	MangaStatusCompleted  MangaStatus = "completed"
	MangaStatusOnHold     MangaStatus = "on_hold"
	MangaStatusDropped    MangaStatus = "dropped"
	MangaStatusPlanToRead MangaStatus = "plan_to_read"
)

func (s MangaStatus) mangaListApply(v *url.Values)               { v.Set("status", string(s)) }
func (s MangaStatus) updateMyMangaListStatusApply(v *url.Values) { v.Set("status", string(s)) }

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
// user's list with values 0-5.
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
