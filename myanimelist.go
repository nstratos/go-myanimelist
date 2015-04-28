package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AnimeList struct {
	MyInfo MyAnimeInfo `xml:"myinfo"`
	Anime  []Anime     `xml:"anime"`
}

type MangaList struct {
	MyInfo MyMangaInfo `xml:"myinfo"`
	Manga  []Manga     `xml:"manga"`
}

type MyAnimeInfo struct {
	MyInfo
	Watching    int `xml:"user_watching"`
	PlanToWatch int `xml:"user_plantowatch"`
}

type MyMangaInfo struct {
	MyInfo
	Reading    int `xml:"user_reading"`
	PlanToRead int `xml:"user_plantoread"`
}

type MyInfo struct {
	ID                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
}

type Anime struct {
	Series
	SeriesAnimedbID   int `xml:"series_animedb_id"`
	SeriesEpisodes    int `xml:"series_episodes"`
	MyWatchedEpisodes int `xml:"my_watched_episodes"`
	My
}

type Manga struct {
	Series
	SeriesMangadbID int `xml:"series_mangadb_id"`
	SeriesChapters  int `xml:"series_chapters"`
	SeriesVolumes   int `xml:"series_volumes"`
	MyReadChapters  int `xml:"my_read_chapters"`
	MyReadVolumes   int `xml:"my_read_volumes"`
	My
}

type Series struct {
	SeriesTitle    string `xml:"series_title"`
	SeriesSynonyms string `xml:"series_synonyms"`
	SeriesType     int    `xml:"series_type"`
	SeriesStatus   int    `xml:"series_status"`
	SeriesStart    string `xml:"series_start"`
	SeriesEnd      string `xml:"series_end"`
	SeriesImage    string `xml:"series_image"`
}

// MyStatus: 1 = watching, 2 = completed, 3 = on hold, 4 = dropped, 6 = plantowatch
type My struct {
	MyId           int    `xml:"my_id"`
	MyStartDate    string `xml:"my_start_date"`
	MyFinishDate   string `xml:"my_finish_date"`
	MyScore        int    `xml:"my_score"`
	MyStatus       int    `xml:"my_status"`
	MyRewatching   string `xml:"my_rewatching"`
	MyRewatchingEp int    `xml:"my_rewatching_ep"`
	MyLastUpdated  string `xml:"my_last_updated"`
	MyTags         string `xml:"my_tags"`
}

func UserAnimeList(username string) (AnimeList, error) {
	const animeListURL = "http://myanimelist.net/malappinfo.php?status=all&type=anime&u="
	//xmlData := readXml()
	xmlData, err := getList(animeListURL + username)
	if err != nil {
		return AnimeList{}, fmt.Errorf("get anime list failed: %s", err)
	}

	al := AnimeList{}
	err = xml.Unmarshal(xmlData, &al)
	if err != nil {
		return AnimeList{}, fmt.Errorf("cannot unmarshal '%s' (%s)", string(xmlData), err)
	}
	return al, nil
}

func UserMangaList(username string) (MangaList, error) {
	const mangaListURL = "http://myanimelist.net/malappinfo.php?status=all&type=manga&u="
	xmlData, err := getList(mangaListURL + username)
	if err != nil {
		return MangaList{}, fmt.Errorf("get manga list failed: %s", err)
	}

	ml := MangaList{}
	err = xml.Unmarshal(xmlData, &ml)
	if err != nil {
		return MangaList{}, fmt.Errorf("cannot unmarshal '%s' (%s)", string(xmlData), err)
	}
	return ml, nil
}

// Fetches a fresh XML from MyAnimeList.net
func getList(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Reads xml from locally stored file myanimelist.xml
func readXml() []byte {
	xmlFile, err := os.Open("./myanimelist.xml")
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
