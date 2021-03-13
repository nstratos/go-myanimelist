package mal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUserServiceMyInfo(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/@me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	u, _, err := client.User.MyInfo(ctx)
	if err != nil {
		t.Errorf("User.MyInfo returned error: %v", err)
	}
	want := &User{ID: 1}
	if got := u; !reflect.DeepEqual(got, want) {
		t.Errorf("User.MyInfo returned\nhave: %+v\n\nwant: %+v", got, want)
	}
}

func TestUserServiceMyInfoError(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/@me", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		http.Error(w, `{"message":"","error":"not_found"}`, 404)
	})

	ctx := context.Background()
	_, _, err := client.User.MyInfo(ctx)
	if err == nil {
		t.Fatal("User.MyInfo expected not found error, got no error.")
	}
	testErrorResponse(t, err, ErrorResponse{Err: "not_found"})
}
