package mal_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

func newOAuth2Client(ctx context.Context) *http.Client {
	// In order to create a client ID and secret for your application:
	//
	//  1. Navigate to https://myanimelist.net/apiconfig or go to your MyAnimeList
	//     profile, click Edit Profile and select the API tab on the far right.
	//  2. Click Create ID and submit the form with your application details.
	oauth2Conf := &oauth2.Config{
		ClientID:     "<Enter your MyAnimeList.net application client ID>",
		ClientSecret: "<Enter your MyAnimeList.net application client secret>", // Optional if you chose App Type 'other'.
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	// To get your first token you need to complete the oauth2 flow. There is a
	// detailed example that uses the terminal under `example/malauth` which you
	// should adjust for your application.
	//
	// Here we assume we have already received our first valid token and stored
	// it somewhere in JSON format.
	const storedToken = `
	{
		"token_type": "Bearer",
		"access_token": "yourAccessToken",
		"refresh_token": "yourRefreshToken",
		"expiry": "2021-06-01T16:12:56.1319122Z"
	}`

	// Decode stored token to oauth2.Token struct.
	oauth2Token := new(oauth2.Token)
	_ = json.Unmarshal([]byte(storedToken), oauth2Token)

	// The oauth2 client returned here with the above configuration and valid
	// token will refresh the token seamlessly when it expires.
	return oauth2Conf.Client(ctx, oauth2Token)
}

func Example_oAuth2() {
	ctx := context.Background()
	oauth2Client := newOAuth2Client(ctx)

	c := mal.NewClient(oauth2Client)

	user, _, err := c.User.MyInfo(ctx)
	if err != nil {
		fmt.Printf("User.MyInfo error: %v", err)
		return
	}
	fmt.Printf("ID: %5d, Joined: %v, Username: %s\n", user.ID, user.JoinedAt.Format("Jan 2006"), user.Name)
}
