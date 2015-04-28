package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AnimeResult struct {
	Entries []AnimeEntry `xml:"entry"`
}

type MangaResult struct {
	Entries []MangaEntry `xml:"entry"`
}

type Entry struct {
	ID        int     `xml:"id"`
	Title     string  `xml:"title"`
	English   string  `xml:"english"`
	Synonyms  string  `xml:"synonyms"`
	Score     float64 `xml:"score"`
	Type      string  `xml:"type"`
	Status    string  `xml:"status"`
	StartDate string  `xml:"start_date"`
	EndDate   string  `xml:"end_date"`
	Synopsis  string  `xml:"synopsis"`
	Image     string  `xml:"image"`
}

type AnimeEntry struct {
	Entry
	Episodes int `xml:"episodes"`
}

type MangaEntry struct {
	Entry
	Chapters int `xml:"chapters"`
	Volumes  int `xml:"volumes"`
}

func SearchAnime(query string) (AnimeResult, error) {
	const searchAnimeURL = "http://myanimelist.net/api/anime/search.xml?q="
	xmlData, err := search(searchAnimeURL, query)
	if err != nil {
		return AnimeResult{}, fmt.Errorf("anime search failed: %s", err)
	}

	result := AnimeResult{}
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return AnimeResult{}, fmt.Errorf("cannot unmarshal '%s' (%s)", string(xmlData), err)
	}

	return result, nil
}

func SearchManga(query string) (MangaResult, error) {
	const searchMangaURL = "http://myanimelist.net/api/manga/search.xml?q="
	xmlData, err := search(searchMangaURL, query)
	if err != nil {
		return MangaResult{}, fmt.Errorf("manga search failed: %s", err)
	}

	result := MangaResult{}
	err = xml.Unmarshal(xmlData, &result)
	if err != nil {
		return MangaResult{}, fmt.Errorf("cannot unmarshal '%s' (%s)", string(xmlData), err)
	}

	return result, nil
}

func search(searchURL, query string) ([]byte, error) {
	req, err := http.NewRequest("GET", searchURL+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.SetBasicAuth(username, password)

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
