package mal_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

func ExampleUserService_MyInfo() {
	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

	// Use a stub server instead of the real API.
	server := newStubServer()
	defer server.Close()
	c.BaseURL, _ = url.Parse(server.URL)

	user, _, err := c.User.MyInfo(ctx)
	if err != nil {
		fmt.Printf("User.MyInfo error: %v", err)
		return
	}
	fmt.Printf("ID: %5d, Joined: %v, Username: %s\n", user.ID, user.JoinedAt.Format("Jan 2006"), user.Name)
	// Output:
	// ID: 4592783, Joined: May 2015, Username: nstratos
}
