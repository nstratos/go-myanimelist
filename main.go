package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const getListUrl = "http://myanimelist.net/malappinfo.php"

const userAgent = `Mozilla/5.0 (Windows NT 6.3; Win64; x64) 
				   AppleWebKit/537.36 (KHTML, like Gecko) 
				   Chrome/37.0.2049.0 Safari/537.36`

type MyAnimeList struct {
	Info  MyInfo  `xml:"myinfo"`
	Anime []Anime `xml:"anime"`
}

type MyInfo struct {
	Id                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Watching          int    `xml:"user_watching"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	PlanToWatch       int    `xml:"user_plantowatch"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
}

// MyStatus: 1 = watching, 2 = completed, 3 = on hold, 4 = dropped, 6 = plantowatch
type Anime struct {
	SeriesAnimedbId   int    `xml:"series_animedb_id"`
	SeriesTitle       string `xml:"series_title"`
	SeriesSynonyms    string `xml:"series_synonyms"`
	SeriesType        int    `xml:"series_type"`
	SeriesEpisodes    int    `xml:"series_episodes"`
	SeriesStatus      int    `xml:"series_status"`
	SeriesStart       string `xml:"series_start"`
	SeriesEnd         string `xml:"series_end"`
	SeriesImage       string `xml:"series_image"`
	MyId              int    `xml:"my_id"`
	MyWatchedEpisodes int    `xml:"my_watched_episodes"`
	MyStartDate       string `xml:"my_start_date"`
	MyFinishDate      string `xml:"my_finish_date"`
	MyScore           int    `xml:"my_score"`
	MyStatus          int    `xml:"my_status"`
	MyRewatching      string `xml:"my_rewatching"`
	MyRewatchingEp    int    `xml:"my_rewatching_ep"`
	MyLastUpdated     string `xml:"my_last_updated"`
	MyTags            string `xml:"my_tags"`
}

func main() {

	xmlData := readXml()

	mal := MyAnimeList{}

	if err := xml.Unmarshal(xmlData, &mal); err != nil {
		log.Fatalln("Error unmarshaling:", err)
		return
	}

	// fmt.Printf("%v \n", mal)
	for _, anime := range mal.Anime {
		fmt.Printf("%s\n", anime.MyRewatching)
	}

}

// Reads xml from locally stored file myanimelist.xml
func readXml() []byte {
	xmlFile, err := os.Open("myanimelist.xml")
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}
	defer xmlFile.Close()
	xml, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	return xml
}

// Fetches a fresh XML from MyAnimeList.net
func fetchXml() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", getListUrl+"?u=Leonteus&status=all&type=anime", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}
