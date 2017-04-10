/*
Package mal provides a client for accessing the MyAnimeList API:
https://myanimelist.net/modules.php?go=api.

Installation

This package can be installed using:

	go get github.com/nstratos/go-myanimelist/mal

Usage

Import the package using:

	import "github.com/nstratos/go-myanimelist/mal"

First construct a new mal client:

	c := mal.NewClient(nil)

Then use one of the client's services (Account, Anime or Manga) to access the
different MyAnimeList API methods.

For example, to get the anime and manga list of the user "Xinil":

	c := mal.NewClient(nil)

	list, _, err := c.Anime.List("Xinil")
	// ...

	list, _, err := c.Manga.List("Xinil")
	// ...


If a method requires authentication, make sure to set your MyAnimeList username
and password on the client.

For example to search for anime and manga (needs authentication):

	c := mal.NewClient(nil)
	c.SetCredentials("<your username>", "<your password>")

	result, _, err := c.Anime.Search("bebop")
	// ...

	result, _, err := c.Manga.Search("bebop")
	// ...

For more complex searches, you can provide the % operator which is escaped as
%% in Go. Note: This is an undocumented API feature.

	c := mal.NewClient(nil)
	c.SetCredentials("<your username>", "<your password>")

	result, _, err := c.Anime.Search("fate%%heaven%%flower")
	// ...
	// Will return: Fate/stay night Movie: Heaven's Feel - I. presage flower

If you need more control, when creating a new client you can pass an
http.Client as an argument.

For example this http.Client passed to the mal client will make sure to cancel
any request that takes longer than 1 second:

	httpcl := &http.Client{
		Timeout: 1 * time.Second,
	}
	c := mal.NewClient(httpcl)
	// ...

See more examples: https://godoc.org/github.com/nstratos/go-myanimelist/mal#pkg-examples

Unit Testing

To run all unit tests:

	cd $GOPATH/src/github.com/nstratos/go-myanimelist/mal
	go test -cover

To see test coverage in your browser:

	go test -covermode=count -coverprofile=count.out && go tool cover -html count.out

Integration Testing

The integration tests will exercise the entire package against the live
MyAnimeList API. As a result, these tests take much longer to run and there is
also a much higher chance of false positives in test failures due to network
issues etc.

These tests are meant to be run using a dedicated test account that contains
empty anime and manga lists. The username and password of the test account need
to be provided every time.

To run the integration tests:

	cd $GOPATH/src/github.com/nstratos/go-myanimelist/mal
	go test -tags=integration -username '<test account username>' -password '<test account password>'

License

MIT

*/
package mal
