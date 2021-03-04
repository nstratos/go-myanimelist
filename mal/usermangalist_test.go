package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestUserServiceMangaList(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/foo/mangalist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{
			"status": "completed",
			"sort":   "manga_id",
			"fields": "foo,bar",
			"limit":  "10",
			"offset": "0",
		})
		const out = `
		{
		  "data": [
		    {
		      "node": { "id": 1 },
			  "list_status": {
			    "status": "plan_to_read"
			  }
		    },
		    {
		      "node": { "id": 2 },
			  "list_status": {
			    "status": "reading"
			  }
		    }
		  ],
		  "paging": {
		    "next": "?offset=4",
		    "previous": "?offset=2"
		  }
		}`
		fmt.Fprintf(w, out)
	})

	ctx := context.Background()
	got, resp, err := client.User.MangaList(ctx, "foo",
		MangaStatusCompleted,
		SortMangaListByMangaID,
		Fields{"foo", "bar"},
		Limit(10),
		Offset(0),
	)
	if err != nil {
		t.Errorf("User.MangaList returned error: %v", err)
	}
	want := []UserManga{
		{
			Manga:  Manga{ID: 1},
			Status: MangaListStatus{Status: "plan_to_read"},
		},
		{
			Manga:  Manga{ID: 2},
			Status: MangaListStatus{Status: "reading"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.MangaList returned\nhave: %+v\n\nwant: %+v", got, want)
	}
	testResponseOffset(t, resp, 4, 2, "User.MangaList")
}

func TestUserServiceMangaListError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/foo/mangalist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"mal is down","error":"internal"}`, 500)
	})

	ctx := context.Background()
	_, resp, err := client.User.MangaList(ctx, "foo")
	if err == nil {
		t.Fatal("User.MangaList expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusInternalServerError, "User.MangaList")
	testErrorResponse(t, err, ErrorResponse{Message: "mal is down", Err: "internal"})
}
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

func TestMangaServiceDeleteMyListItem(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/manga/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	ctx := context.Background()
	resp, err := client.Manga.DeleteMyListItem(ctx, 1)
	if err != nil {
		t.Errorf("Manga.DeleteMyListItem returned error: %v", err)
	}
	testResponseStatusCode(t, resp, http.StatusOK, "Manga.DeleteMyListItem")
}

func TestMangaServiceDeleteMyListItemError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/manga/1/my_list_status", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		http.Error(w, `{"message":"manga not found","error":"not_found"}`, http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Manga.DeleteMyListItem(ctx, 1)
	if err == nil {
		t.Fatal("Manga.DeleteMyListItem expected internal error, got no error.")
	}
	testResponseStatusCode(t, resp, http.StatusNotFound, "Manga.DeleteMyListItem")
	testErrorResponse(t, err, ErrorResponse{Message: "manga not found", Err: "not_found"})
}
