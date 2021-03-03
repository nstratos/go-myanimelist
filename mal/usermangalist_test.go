package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestMangaServiceUpdateMyListStatus(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/manga/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		const out = `
		{
		  "status": "completed",
		  "score": 8,
		  "num_volumes_read": 3,
		  "num_chapters_read": 3,
		  "is_rereading": true,
		  "updated_at": "2018-04-25T15:59:52Z",
		  "priority": 2,
		  "num_times_reread": 2,
		  "reread_value": 1,
		  "tags": ["foo","bar"],
		  "comments": "comments"
		}`
		fmt.Fprintf(w, out)
	})

	ctx := context.Background()
	got, _, err := client.Manga.UpdateMyListStatus(ctx, 1,
		MangaStatusCompleted,
		IsRereading(true),
		Score(8),
		NumVolumesRead(3),
		NumChaptersRead(3),
		Priority(2),
		NumTimesReread(2),
		RereadValue(1),
		Tags{"foo", "bar"},
		Comments("comments"),
	)
	if err != nil {
		t.Errorf("Manga.UpdateMyListStatus returned error: %v", err)
	}

	want := &MangaListStatus{
		Status:          MangaStatusCompleted,
		IsRereading:     true,
		Score:           8,
		NumVolumesRead:  3,
		NumChaptersRead: 3,
		Priority:        2,
		NumTimesReread:  2,
		RereadValue:     1,
		Tags:            []string{"foo", "bar"},
		Comments:        "comments",
		UpdatedAt:       time.Date(2018, 04, 25, 15, 59, 52, 0, time.UTC),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Manga.UpdateMyListStatus returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestMangaServiceUpdateMyListStatusError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/manga/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		http.Error(w, `{"message":"mal is down","error":"internal"}`, 500)
	})

	ctx := context.Background()
	_, resp, err := client.Manga.UpdateMyListStatus(ctx, 1)
	if err == nil {
		t.Fatal("Manga.UpdateMyListStatus expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusInternalServerError, "Manga.UpdateMyListStatus")
	testErrorResponse(t, err, ErrorResponse{Message: "mal is down", Err: "internal"})
}
