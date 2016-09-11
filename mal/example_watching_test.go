package mal_test

import (
	"fmt"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
)

func Example_printWatchingAnime() {
	list, _, err := mal.NewClient(nil).Anime.List("Xinil")
	if err != nil {
		log.Fatal(err)
	}
	for _, anime := range list.Anime {
		// If anime has status watching then print it.
		if anime.MyStatus == mal.StatusWatching {
			fmt.Printf("%s\n", anime.SeriesTitle)
		}
	}
}
