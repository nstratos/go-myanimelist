package main

import (
	"fmt"
	"log"

	"bitbucket.org/nstratos/mal"
)

func main() {

	//getAnime("Leonteus")
	//search("Full metal")
	verify("Leonteus", "001010100")
}

func verify(username, password string) {
	user, err := mal.Verify(username, password)
	if err != nil {
		log.Fatalf("Error verifying: %s\n", err)
	}
	fmt.Printf("%+v\n", user)
}

func getAnime(username string) {
	myAnimeList, err := mal.GetAnime(username)
	if err != nil {
		log.Fatalf("Error getting list: %s\n", err)
		return
	}
	for _, anime := range myAnimeList.Anime {
		fmt.Printf("%s\n", anime.SeriesTitle)
	}
}

func search(query string) {
	result, err := mal.Search(query)
	if err != nil {
		log.Fatalf("Error searching: %s\n", err)
		return
	}
	printResults(result.Entries)
}

func printResults(entries []mal.Entry) {
	for _, entry := range entries {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("| Title: %s\n", entry.Title)
		fmt.Printf("| Episodes: %d\n", entry.Episodes)
		fmt.Printf("| Type: %s\n", entry.Type)
		fmt.Printf("| Synopsis: %s\n", entry.Synopsis)
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("\n")
	}
}
