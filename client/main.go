package main

import (
	"fmt"
	"log"

	"bitbucket.org/nstratos/mal"
)

/*const userAgent = `Mozilla/5.0 (Windows NT 6.3; Win64; x64)
  AppleWebKit/537.36 (KHTML, like Gecko)
  Chrome/37.0.2049.0 Safari/537.36`*/

func main() {

	//mal.Init("Leonteus", "001010100", "api-indiv-2D4068FCF43349DA30D8D4E5667883C2")
	//verify()
	//searchManga("naruto")
	//getAnime("Leonteus")
	//getManga("Leonteus")
	//data := mal.AnimeData{Status: "1", Score: 9}
	// data := mal.AnimeData{
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
	client := mal.NewClient()
	client.SetCredentials("Leonteus", "001010100")
	client.SetUserAgent("api-indiv-2D4068FCF43349DA30D8D4E5667883C2")
	// _, err := client.Anime.Delete(11933)
	// if err != nil {
	// 	log.Fatalf("Anime.Delete error: %v\n", err)
	// }
	_, err := client.Anime.Add(11933, mal.AnimeData{Status: "1", Score: 10})
	if err != nil {
		log.Fatalf("Anime.Add error: %v\n", err)
	}
}

func verify() {
	user, err := mal.Verify()
	if err != nil {
		log.Fatalf("Error verifying: %s\n", err)
	}
	fmt.Printf("%+v\n", user)
}

func getAnime(username string) {
	al, err := mal.UserAnimeList(username)
	if err != nil {
		log.Fatalf("getAnime: %s\n", err)
		return
	}

	for _, anime := range al.Anime {
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
		fmt.Printf("| SeriesAnimedbID: %v\n", anime.SeriesAnimedbID)
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
	fmt.Printf("| ID: %v\n", al.MyInfo.ID)
	fmt.Printf("| Name: %v\n", al.MyInfo.Name)
	fmt.Printf("| Completed: %v\n", al.MyInfo.Completed)
	fmt.Printf("| OnHold: %v\n", al.MyInfo.OnHold)
	fmt.Printf("| Dropped: %v\n", al.MyInfo.Dropped)
	fmt.Printf("| DaysSpentWatching: %v\n", al.MyInfo.DaysSpentWatching)
	fmt.Printf("| Watching: %v\n", al.MyInfo.Watching)
	fmt.Printf("| PlanToWatch: %v\n", al.MyInfo.PlanToWatch)

}

func getManga(username string) {
	ml, err := mal.UserMangaList(username)
	if err != nil {
		log.Fatalf("getManga: %s\n", err)
		return
	}

	for _, manga := range ml.Manga {
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
		fmt.Printf("| SeriesMangadbID: %v\n", manga.SeriesMangadbID)
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
	fmt.Printf("| ID: %v\n", ml.MyInfo.ID)
	fmt.Printf("| Name: %v\n", ml.MyInfo.Name)
	fmt.Printf("| Completed: %v\n", ml.MyInfo.Completed)
	fmt.Printf("| OnHold: %v\n", ml.MyInfo.OnHold)
	fmt.Printf("| Dropped: %v\n", ml.MyInfo.Dropped)
	fmt.Printf("| DaysSpentWatching: %v\n", ml.MyInfo.DaysSpentWatching)
	fmt.Printf("| Reading: %v\n", ml.MyInfo.Reading)
	fmt.Printf("| PlanToRead: %v\n", ml.MyInfo.PlanToRead)

}

func searchAnime(query string) {
	result, err := mal.SearchAnime(query)
	if err != nil {
		log.Fatalf("searchAnime: %s\n", err)
		return
	}
	for _, entry := range result.Entries {
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

func searchManga(query string) {
	result, err := mal.SearchManga(query)
	if err != nil {
		log.Fatalf("Error searching: %s\n", err)
		return
	}
	for _, entry := range result.Entries {
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
