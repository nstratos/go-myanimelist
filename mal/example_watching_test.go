package mal_test

import (
	"fmt"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
)

func Example_printWatchingAnime() {
	list, _, err := mal.NewClient().Anime.List("Xinil")
	if err != nil {
		log.Fatal(err)
	}
	for _, anime := range list.Anime {
		// If anime has its status set as currently watching then print it.
		if anime.MyStatus == mal.Current {
			fmt.Printf("%s\n", anime.SeriesTitle)
		}
	}
}
