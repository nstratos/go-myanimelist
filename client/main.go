package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"bitbucket.org/nstratos/mal"
)

/*const userAgent = `Mozilla/5.0 (Windows NT 6.3; Win64; x64)
  AppleWebKit/537.36 (KHTML, like Gecko)
  Chrome/37.0.2049.0 Safari/537.36`*/

func main() {

	//verify()
	//searchManga("naruto")
	//getAnime("Leonteus")
	//getManga("Leonteus")
	//data := mal.AnimeEntry{Status: "1", Score: 9}
	// data := mal.AnimeEntry{
	// 	Episode:            1,
	// 	Status:             "onhold",
	// 	Score:              4,
	// 	DownloadedEpisodes: 1,
	// 	StorageType:        0,
	// 	StorageValue:       0.0,
	// 	TimesRewatched:     1,
	// 	RewatchValue:       5,
	// 	DateStart:          "12252007",
	// 	DateFinish:         "12252008",
	// 	Priority:           3,
	// 	EnableDiscussion:   0,
	// 	EnableRewatching:   1,
	// 	Comments:           "good",
	// 	FansubGroup:        "horriblesubs",
	// 	Tags:               "mytag",
	// }
	//mal.UpdateAnime(9989, data)
	//mal.AddAnime(9989, data)
	//searchAnime("anohana")
	//mal.DeleteAnime(9989)
	//searchManga("anohana")
	//data := mal.MangaData{Status: "1", Score: 10}
	//mal.AddManga(35733, data)
	//mal.UpdateManga(35733, data)
	//mal.DeleteManga(35733)
	agent, err := ioutil.ReadFile("agent.txt")
	if err != nil {
		log.Fatalln("cannot read agent:", err)
	}
	client := mal.NewClient()
	client.SetCredentials("Leonteus", "001010100")
	client.SetUserAgent(string(agent))
	// _, err := client.Anime.Delete(11933)
	// if err != nil {
	// 	log.Fatalf("Anime.Delete error: %v\n", err)
	// }
	// _, err := client.Anime.Add(11933, mal.AnimeEntry{Status: "1", Score: 10})
	// if err != nil {
	// 	log.Fatalf("Anime.Add error: %v\n", err)
	// }
	result, resp, err := client.Manga.Search("anohana")
	if err != nil {
		log.Fatalf("Anime.Search error: %v, received: '%v'\n", err, string(resp.Body))
		//log.Fatalf("Anime.Search error: %v\n", err)
	}
	printMangaResult(result)
}

func printAnimeList(list mal.AnimeList) {
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
}

func printMangaList(list mal.MangaList) {
	for _, manga := range list.Manga {
		fmt.Printf("----------------Manga-------------------\n")
		fmt.Printf("| MyID: %v\n", manga.MyID)
		fmt.Printf("| MyStartDate: %v\n", manga.MyStartDate)
		fmt.Printf("| MyFinishDate: %v\n", manga.MyFinishDate)
		fmt.Printf("| MyScore: %v\n", manga.MyScore)
		fmt.Printf("| MyStatus: %v\n", manga.MyStatus)
		fmt.Printf("| MyRewatching: %v\n", manga.MyRewatching)
		fmt.Printf("| MyRewatchingEp: %v\n", manga.MyRewatchingEp)
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
}

func printAnimeResult(result *mal.AnimeResult) {
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

func printMangaResult(result *mal.MangaResult) {
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
