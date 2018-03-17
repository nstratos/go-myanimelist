package mal_test

import (
	"fmt"
	"log"

	"github.com/nstratos/go-myanimelist/mal"
)

// anime examples

func ExampleAnimeService_Add() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	resp, err := c.Anime.Add(9989, mal.AnimeEntry{Status: mal.Current, Episode: 1})
	if err != nil {
		log.Fatalf("Anime.Add error: %v, received: '%v'\n", err, string(resp.Body))
	}
}

func ExampleAnimeService_Update() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	resp, err := c.Anime.Update(9989, mal.AnimeEntry{Status: mal.Completed, Score: 9})
	if err != nil {
		log.Fatalf("Anime.Update error: %v, received: '%v'\n", err, string(resp.Body))
	}
}

func ExampleAnimeService_Search() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	result, resp, err := c.Anime.Search("anohana")
	if err != nil {
		log.Fatalf("Anime.Search error: %v, received: '%v'\n", err, string(resp.Body))
	}

	// For more complex searches, you can provide the % operator which is
	// escaped as %% in Go. Note: This is an undocumented API feature.
	//
	// As an example, if you search for "fate%%heaven%%flower" you can get one
	// accurate result of the title:
	//
	// Fate/stay night Movie: Heaven's Feel - I. presage flower

	// printing results
	for _, entry := range result.Rows {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("| ID: %v\n", entry.ID)
		fmt.Printf("| Title: %v\n", entry.Title)
		fmt.Printf("| Episodes: %d\n", entry.Episodes)
		fmt.Printf("| Type: %v\n", entry.Type)
		fmt.Printf("| Score: %v\n", entry.Score)
		fmt.Printf("| Synopsis: %v\n", entry.Synopsis)
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("\n")
	}
}

func ExampleAnimeService_List() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	list, resp, err := c.Anime.List("Xinil")
	if err != nil {
		log.Fatalf("Anime.List error: %v, received: '%v'\n", err, string(resp.Body))
	}

	// printing results
	for _, anime := range list.Anime {
		fmt.Printf("----------------Anime-------------------\n")
		fmt.Printf("| MyID: %v\n", anime.MyID)
		fmt.Printf("| MyStartDate: %v\n", anime.MyStartDate)
		fmt.Printf("| MyFinishDate: %v\n", anime.MyFinishDate)
		fmt.Printf("| MyScore: %v\n", anime.MyScore)
		fmt.Printf("| MyStatus: %v\n", anime.MyStatus)
		fmt.Printf("| MyRewatching: %v\n", anime.MyRewatching)
		fmt.Printf("| MyRewatchingEp: %v\n", anime.MyRewatchingEp)
		fmt.Printf("| MyLastUpdated: %v\n", anime.MyLastUpdated)
		fmt.Printf("| MyTags: %v\n", anime.MyTags)
		fmt.Printf("| MyWatchedEpisodes: %v\n", anime.MyWatchedEpisodes)
		fmt.Printf("| SeriesAnimeDBID: %v\n", anime.SeriesAnimeDBID)
		fmt.Printf("| SeriesEpisodes: %v\n", anime.SeriesEpisodes)
		fmt.Printf("| SeriesTitle: %v\n", anime.SeriesTitle)
		fmt.Printf("| SeriesSynonyms: %v\n", anime.SeriesSynonyms)
		fmt.Printf("| SeriesType: %v\n", anime.SeriesType)
		fmt.Printf("| SeriesStatus: %v\n", anime.SeriesStatus)
		fmt.Printf("| SeriesStart: %v\n", anime.SeriesStart)
		fmt.Printf("| SeriesEnd: %v\n", anime.SeriesEnd)
		fmt.Printf("| SeriesImage: %v\n", anime.SeriesImage)
		fmt.Printf("\n")
	}
	fmt.Printf("----------------MyInfo------------------\n")
	fmt.Printf("| ID: %v\n", list.MyInfo.ID)
	fmt.Printf("| Name: %v\n", list.MyInfo.Name)
	fmt.Printf("| Completed: %v\n", list.MyInfo.Completed)
	fmt.Printf("| OnHold: %v\n", list.MyInfo.OnHold)
	fmt.Printf("| Dropped: %v\n", list.MyInfo.Dropped)
	fmt.Printf("| DaysSpentWatching: %v\n", list.MyInfo.DaysSpentWatching)
	fmt.Printf("| Watching: %v\n", list.MyInfo.Watching)
	fmt.Printf("| PlanToWatch: %v\n", list.MyInfo.PlanToWatch)
	fmt.Printf("----------------------------------------\n")
}

// manga examples

func ExampleMangaService_Add() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	resp, err := c.Manga.Add(35733, mal.MangaEntry{Status: mal.Current, Chapter: 1, Volume: 1})
	if err != nil {
		log.Fatalf("Manga.Add error: %v, received: '%v'\n", err, string(resp.Body))
	}
}

func ExampleMangaService_Update() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	resp, err := c.Manga.Update(35733, mal.MangaEntry{Status: mal.Completed, Score: 9})
	if err != nil {
		log.Fatalf("Manga.Update error: %v, received: '%v'\n", err, string(resp.Body))
	}
}

func ExampleMangaService_Search() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	result, resp, err := c.Manga.Search("anohana")
	if err != nil {
		log.Fatalf("Manga.Search error: %v, received: '%v'\n", err, string(resp.Body))
	}

	for _, entry := range result.Rows {
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("| ID: %v\n", entry.ID)
		fmt.Printf("| Title: %v\n", entry.Title)
		fmt.Printf("| Chapters: %d\n", entry.Chapters)
		fmt.Printf("| Volumes: %d\n", entry.Volumes)
		fmt.Printf("| Type: %v\n", entry.Type)
		fmt.Printf("| Score: %v\n", entry.Score)
		fmt.Printf("| Synopsis: %v\n", entry.Synopsis)
		fmt.Printf("----------------------------------------\n")
		fmt.Printf("\n")
	}
}

func ExampleMangaService_List() {
	c := mal.NewClient(mal.Auth("<your username>", "<your password>"))

	list, resp, err := c.Manga.List("Xinil")
	if err != nil {
		log.Fatalf("Manga.List error: %v, received: '%v'\n", err, string(resp.Body))
	}

	// printing results
	for _, manga := range list.Manga {
		fmt.Printf("----------------Manga-------------------\n")
		fmt.Printf("| MyID: %v\n", manga.MyID)
		fmt.Printf("| MyStartDate: %v\n", manga.MyStartDate)
		fmt.Printf("| MyFinishDate: %v\n", manga.MyFinishDate)
		fmt.Printf("| MyScore: %v\n", manga.MyScore)
		fmt.Printf("| MyStatus: %v\n", manga.MyStatus)
		fmt.Printf("| MyRereading: %v\n", manga.MyRereading)
		fmt.Printf("| MyRereadingChap: %v\n", manga.MyRereadingChap)
		fmt.Printf("| MyLastUpdated: %v\n", manga.MyLastUpdated)
		fmt.Printf("| MyTags: %v\n", manga.MyTags)
		fmt.Printf("| MyReadChapters: %v\n", manga.MyReadChapters)
		fmt.Printf("| MyReadVolumes: %v\n", manga.MyReadVolumes)
		fmt.Printf("| SeriesMangaDBID: %v\n", manga.SeriesMangaDBID)
		fmt.Printf("| SeriesChapters: %v\n", manga.SeriesChapters)
		fmt.Printf("| SeriesVolumes: %v\n", manga.SeriesVolumes)
		fmt.Printf("| SeriesTitle: %v\n", manga.SeriesTitle)
		fmt.Printf("| SeriesSynonyms: %v\n", manga.SeriesSynonyms)
		fmt.Printf("| SeriesType: %v\n", manga.SeriesType)
		fmt.Printf("| SeriesStatus: %v\n", manga.SeriesStatus)
		fmt.Printf("| SeriesStart: %v\n", manga.SeriesStart)
		fmt.Printf("| SeriesEnd: %v\n", manga.SeriesEnd)
		fmt.Printf("| SeriesImage: %v\n", manga.SeriesImage)
		fmt.Printf("\n")
	}
	fmt.Printf("----------------MyInfo------------------\n")
	fmt.Printf("| ID: %v\n", list.MyInfo.ID)
	fmt.Printf("| Name: %v\n", list.MyInfo.Name)
	fmt.Printf("| Completed: %v\n", list.MyInfo.Completed)
	fmt.Printf("| OnHold: %v\n", list.MyInfo.OnHold)
	fmt.Printf("| Dropped: %v\n", list.MyInfo.Dropped)
	fmt.Printf("| DaysSpentWatching: %v\n", list.MyInfo.DaysSpentWatching)
	fmt.Printf("| Reading: %v\n", list.MyInfo.Reading)
	fmt.Printf("| PlanToRead: %v\n", list.MyInfo.PlanToRead)
	fmt.Printf("----------------------------------------\n")
}
