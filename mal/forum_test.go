package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestForumServiceBoards(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/forum/boards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLValues(t, r, urlValues{})
		testBody(t, r, "")
		fmt.Fprintf(w, `
		{
		  "categories": [
		    {
		  	  "title": "MyAnimeList",
		  	  "boards": [
		        {
		          "id": 17,
		          "title": "MAL Guidelines",
		          "description": "Site rules.",
		          "subboards": [{"id": 2,"title": "Anime DB"}]
		  	    }
		      ]
		    }
		  ]
		}`)
	})

	ctx := context.Background()
	a, _, err := client.Forum.Boards(ctx)
	if err != nil {
		t.Errorf("Forum.Boards returned error: %v", err)
	}
	want := &Forum{
		Categories: []ForumCategory{
			{
				Title: "MyAnimeList",
				Boards: []ForumBoard{
					{
						ID:          17,
						Title:       "MAL Guidelines",
						Description: "Site rules.",
						Subboards:   []ForumSubboard{{ID: 2, Title: "Anime DB"}},
					},
				},
			},
		},
	}
	if got := a; !reflect.DeepEqual(got, want) {
		t.Errorf("Forum.Boards returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestForumServiceBoardsError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/forum/boards", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"forum deleted","error":"not_found"}`, 404)
	})

	ctx := context.Background()
	_, _, err := client.Forum.Boards(ctx)
	if err == nil {
		t.Fatal("Forum.Boards expected not found error, got no error.")
	}
	testErrorResponse(t, err, ErrorResponse{Message: "forum deleted", Err: "not_found"})
}
