# go-myanimelist #

go-myanimelist is a Go library for accessing the [MyAnimeList API](http://myanimelist.net/modules.php?go=api).

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/nstratos/go-myanimelist/mal?status.svg)](https://godoc.org/github.com/nstratos/go-myanimelist/mal)
[![Coverage Status](https://coveralls.io/repos/nstratos/go-myanimelist/badge.svg?branch=master)](https://coveralls.io/r/nstratos/go-myanimelist?branch=master)
[![Build Status](https://drone.io/github.com/nstratos/go-myanimelist/status.png)](https://drone.io/github.com/nstratos/go-myanimelist/latest)

## Installation ##

This package can be installed using:

    go get github.com/nstratos/go-myanimelist/mal

## Usage ##

```go
import "github.com/nstratos/go-myanimelist/mal"
```

Construct a new client, then use one of the client's services to access the
different MyAnimeList API methods. For example, to get the anime list of the
user "Xinil":

```go
c := mal.NewClient(nil)
c.SetCredentials("YOUR_MYANIMELIST_USERNAME", "YOUR_MYANIMELIST_PASSWORD")
c.SetUserAgent("YOUR_WHITELISTED_USER_AGENT")

list, _, err := c.Anime.List("Xinil")
// handle err

// do something with list
```

See more [examples](https://godoc.org/github.com/nstratos/go-myanimelist/mal#pkg-examples).
