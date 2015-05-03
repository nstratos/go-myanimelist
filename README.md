# go-myanimelist #

go-myanimelist is a Go client library for accessing the [MyAnimeList API](http://myanimelist.net/modules.php?go=api).

## Installation ## 

This package can be installed using:

    go get github.com/nstratos/go-myanimelist/mal

## Usage ##

	import "github.com/nstratos/go-myanimelist/mal"

Construct a new client, then use one of the client's services to access the
different MyAnimeList API methods. For example, to get the anime list of the
user "Xinil":

	c := mal.NewClient()
	c.SetCredentials("YOUR_MYANIMELIST_USERNAME", "YOUR_MYANIMELIST_PASSWORD")
	c.SetUserAgent("YOUR_WHITELISTED_USER_AGENT")

	list, _, err := c.Anime.List("Xinil")
