/*
Package mal provides a client for accessing the MyAnimeList API.

Construct a new client, then use one of the client's services to access the
different MyAnimeList API methods. For example, to get the anime list of the
user "Xinil":

	c := mal.NewClient(nil)
	c.SetCredentials("YOUR_MYANIMELIST_USERNAME", "YOUR_MYANIMELIST_PASSWORD")
	c.SetUserAgent("YOUR_WHITELISTED_USER_AGENT")

	list, _, err := c.Anime.List("Xinil")
	// handle err

	// do something with list
*/
package mal
