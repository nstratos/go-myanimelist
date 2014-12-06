package main

import (
	"fmt"
	"log"
	"net/url"

	"bitbucket.org/nstratos/mal"
)

func main() {

	// myAnimeList, err := mal.Get()
	// if err != nil {
	// 	log.Fatalf("Error getting list (%s)\n", err)
	// 	return
	// }
	// for _, anime := range myAnimeList.Anime {
	// 	fmt.Printf("%s\n", anime.SeriesTitle)
	// }
	result, err := mal.Search(url.QueryEscape("full metal"))
	if err != nil {
		log.Fatalf("Error searching (%s)\n", err)
		return
	}
	for _, entry := range result.Entries {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("| Title: %s\n", entry.Title)
		fmt.Printf("| Episodes: %d\n", entry.Episodes)
		fmt.Printf("| Type: %s\n", entry.Type)
		fmt.Printf("| Synopsis: %s\n", entry.Synopsis)
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("\n")
	}

}
