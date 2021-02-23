# go-myanimelist

go-myanimelist is a Go client library for accessing the [MyAnimeList API v2](https://myanimelist.net/apiconfig/references/api/v2).

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/nstratos/go-myanimelist/mal?status.svg)](https://godoc.org/github.com/nstratos/go-myanimelist/mal)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstratos/go-myanimelist)](https://goreportcard.com/report/github.com/nstratos/go-myanimelist)
[![Coverage Status](https://coveralls.io/repos/github/nstratos/go-myanimelist/badge.svg?branch=master)](https://coveralls.io/github/nstratos/go-myanimelist?branch=master)
[![Build Status](https://travis-ci.org/nstratos/go-myanimelist.svg?branch=master)](https://travis-ci.org/nstratos/go-myanimelist)

As of March 2017, this package is featured in
[awesome-go](https://github.com/avelino/awesome-go).

## Installation

This package can be installed using:

	go get github.com/nstratos/go-myanimelist/mal

## Usage

Import the package using:

```go
import "github.com/nstratos/go-myanimelist/mal"
```

First construct a new mal client:

```go
c := mal.NewClient()
```

Then use one of the client's services (Account, Anime or Manga) to access the
different MyAnimeList API methods.

## List

To get the anime and manga list of a user:

```go
c := mal.NewClient()

list, _, err := c.Anime.List("Xinil")
// ...

list, _, err := c.Manga.List("Xinil")
// ...
```

## Authentication

Beyond List, the rest of the methods require authentication so typically you
will use an option to pass username and password to NewClient:

```go
c := mal.NewClient(
	mal.Auth("<your username>", "<your password>"),
)
```

## Search

To search for anime and manga:

```go
c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

result, _, err := c.Anime.Search("bebop")
// ...

result, _, err := c.Manga.Search("bebop")
// ...
```

For more complex searches, you can provide the % operator which acts as a
wildcard and is escaped as %% in Go:

```go
result, _, err := c.Anime.Search("fate%%heaven%%flower")
// ...
// Will return: Fate/stay night Movie: Heaven's Feel - I. presage flower
```

Note: This is an undocumented feature of the MyAnimeList Search method.

## Add

To add anime and manga, you provide their IDs and values through AnimeEntry and
MangaEntry:

```go
c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

_, err := c.Anime.Add(9989, mal.AnimeEntry{Status: mal.Current, Episode: 1})
// ...

_, err := c.Manga.Add(35733, mal.MangaEntry{Status: mal.Planned, Chapter: 1, Volume: 1})
// ...
```

Note that when adding entries, Status is required.

## Update

Similar to Add, Update also needs the ID of the entry and the values to be
updated:

```go
c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

_, err := c.Anime.Update(9989, mal.AnimeEntry{Status: mal.Completed, Score: 9})
// ...

_, err := c.Manga.Update(35733, mal.MangaEntry{Status: mal.OnHold})
// ...
```

## Delete

To delete anime and manga, simply provide their IDs:

```go
c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

_, err := c.Anime.Delete(9989)
// ...

_, err := c.Manga.Delete(35733)
// ...
```

## More Examples

See package examples:
https://godoc.org/github.com/nstratos/go-myanimelist/mal#pkg-examples

## Advanced Control

If you need more control over the created requests, you can use an option to
pass a custom HTTP client to NewClient:

```go
c := mal.NewClient(
	mal.HTTPClient(&http.Client{}),
)
```

For example this http.Client will make sure to cancel any request that takes
longer than 1 second:

```go
httpcl := &http.Client{
	Timeout: 1 * time.Second,
}
c := mal.NewClient(mal.HTTPClient(httpcl))
// ...
```

## Unit Testing

To run all unit tests:

	go test -cover

To see test coverage in your browser:

	go test -covermode=count -coverprofile=count.out && go tool cover -html count.out

## Integration Testing

The integration tests will exercise the entire package against the live
MyAnimeList API. As a result, these tests take much longer to run and there is
also a much higher chance of false positives in test failures due to network
issues etc.

These tests are meant to be run using a dedicated test account that contains
empty anime and manga lists. The username and password of the test account need
to be provided every time.

To run the integration tests:

	go test -tags=integration -username '<test account username>' -password '<test account password>'

## License

MIT
