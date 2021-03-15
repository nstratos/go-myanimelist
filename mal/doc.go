/*
Package mal provides a client for accessing the MyAnimeList API:
https://myanimelist.net/apiconfig/references/api/v2.

Installation

This package can be installed using:

	go get github.com/nstratos/go-myanimelist/mal

Usage

Import the package using:

	import "github.com/nstratos/go-myanimelist/mal"

First construct a new mal client:

	c := mal.NewClient(nil)

Then use one of the client's services (User, Anime, Manga and Forum) to access
the different MyAnimeList API methods.

Authentication

When creating a new client, pass an http.Client that can handle authentication
for you. The recommended way is to use the golang.org/x/oauth2 package
(https://github.com/golang/oauth2). After performing the OAuth2 flow, you will
get an access token which can be used like this:

	ctx := context.Background()
	c := mal.NewClient(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "<your access token>"},
		)),
	)

Note that all calls made by the client above will include the specified access
token which is specific for an authenticated user. Therefore, authenticated
clients should almost never be shared between different users.

Performing the OAuth2 flow involves registering a MAL API application and then
asking for the user's consent to allow the application to access their data.

There is a detailed example of how to perform the Oauth2 flow and get an access
token through the terminal under example/malauth. The only thing you need to run
the example is a client ID and a client secret which you can acquire after
registering your MAL API application. Here's how:

 1. Navigate to https://myanimelist.net/apiconfig or go to your MyAnimeList
    profile, click Edit Profile and select the API tab on the far right.

 2. Click Create ID and submit the form with your application details.

After registering your application, you can run the example and pass the client
ID and client secret through flags:

	cd example/malauth
	go run main.go democlient.go --client-id=... --client-secret=...

	or

	go install github.com/nstratos/go-myanimelist/example/malauth
	malauth --client-id=... --client-secret=...

After you perform a successful authentication once, the access token will be
cached in a file under the same directory which makes it easier to run the
example multiple times.

Official MAL API OAuth2 docs:
https://myanimelist.net/apiconfig/references/authorization

List

To search and get anime and manga data:

	list, _, err := c.Anime.List(ctx, "hokuto no ken",
		mal.Fields{"rank", "popularity", "my_list_status"},
		mal.Limit(5),
	)
	// ...

	list, _, err := c.Manga.List(ctx, "hokuto no ken",
		mal.Fields{"rank", "popularity", "my_list_status"},
		mal.Limit(5),
	)
	// ...

You may get user specific data for a certain record by using the optional field
"my_list_status".

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/anime_get

- https://myanimelist.net/apiconfig/references/api/v2#operation/manga_get

UserList

To get the anime or manga list of a user:

	// Get the authenticated user's anime list, filter only watching anime, sort by
	// last updated, include list status.
	anime, _, err := c.User.AnimeList(ctx, "@me",
	    mal.Fields{"list_status"},
	    mal.AnimeStatusWatching,
	    mal.SortAnimeListByListUpdatedAt,
	    mal.Limit(5),
	)
	// ...

	// Get the authenticated user's manga list's second page, sort by score,
	// include list status, comments and tags.
	manga, _, err := c.User.MangaList(ctx, "@me",
	    mal.SortMangaListByListScore,
	    mal.Fields{"list_status{comments, tags}"},
	    mal.Limit(5),
	    mal.Offset(1),
	)
	// ...

You may provide the username of the user or "@me" to get the authenticated
user's list.

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/users_user_id_animelist_get

- https://myanimelist.net/apiconfig/references/api/v2#operation/users_user_id_mangalist_get

MyInfo

To get information about the authenticated user:

	user, _, err := c.User.MyInfo(ctx)
	// ...

This method can use the Fields option but the API doesn't seem to be able to
send optional fields like "anime_statistics" at the time of writing.

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/users_user_id_get

Details

To get details for a certain anime or manga:

	a, _, err := c.Anime.Details(ctx, 967,
		mal.Fields{
			"alternative_titles",
			"media_type",
			"num_episodes",
			"start_season",
			"source",
			"genres",
			"studios",
			"average_episode_duration",
		},
	)
	// ...

	m, _, err := c.Manga.Details(ctx, 401,
		mal.Fields{
			"alternative_titles",
			"media_type",
			"num_volumes",
			"num_chapters",
			"authors{last_name, first_name}",
			"genres",
			"status",
		},
	)
	// ...

By default most fields are not populated so use the Fields option to request the
fields you need.

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/anime_anime_id_get

- https://myanimelist.net/apiconfig/references/api/v2#operation/manga_manga_id_get

Ranking

To get anime or manga based on a certain ranking:

	anime, _, err := c.Anime.Ranking(ctx,
		mal.AnimeRankingAiring,
		mal.Fields{"rank", "popularity"},
		mal.Limit(6),
	)
	// ...

	manga, _, err := c.Manga.Ranking(ctx,
		mal.MangaRankingByPopularity,
		mal.Fields{"rank", "popularity"},
		mal.Limit(6),
	)
	// ...

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/anime_ranking_get

- https://myanimelist.net/apiconfig/references/api/v2#operation/manga_ranking_get

Add

To add anime and manga, you provide their IDs and values through AnimeEntry and
MangaEntry:

	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	_, err := c.Anime.Add(9989, mal.AnimeEntry{Status: mal.Current, Episode: 1})
	// ...

	_, err := c.Manga.Add(35733, mal.MangaEntry{Status: mal.Planned, Chapter: 1, Volume: 1})
	// ...

Note that when adding entries, Status is required.

Update

Similar to Add, Update also needs the ID of the entry and the values to be
updated:

	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	_, err := c.Anime.Update(9989, mal.AnimeEntry{Status: mal.Completed, Score: 9})
	// ...

	_, err := c.Manga.Update(35733, mal.MangaEntry{Status: mal.OnHold})
	// ...

Delete

To delete anime or manga from a user's list, simply provide their IDs:

	_, err := c.Anime.DeleteMyListItem(ctx, 967)
	// ...

	_, err := c.Manga.DeleteMyListItem(ctx, 401)
	// ...

Official docs:

- https://myanimelist.net/apiconfig/references/api/v2#operation/anime_anime_id_my_list_status_delete

- https://myanimelist.net/apiconfig/references/api/v2#operation/manga_manga_id_my_list_status_delete

More Examples

See package examples:
https://pkg.go.dev/github.com/nstratos/go-myanimelist/mal#pkg-examples

Unit Testing

To run all unit tests:

	go test -cover

To see test coverage in your browser:

	go test -covermode=count -coverprofile=count.out && go tool cover -html count.out

Integration Testing

The integration tests will exercise the entire package against the live
MyAnimeList API. As a result, these tests take much longer to run and there is
also a much higher chance of false positives in test failures due to network
issues etc.

These tests are meant to be run using a dedicated test account that contains
empty anime and manga lists. A valid access token needs to be provided every
time. Check the authentication section to learn how to get one.

By default the integration tests are skipped when an access token is not
provided. To run all tests including the integration tests:

	go test --access-token '<your access token>'

License

MIT

*/
package mal
