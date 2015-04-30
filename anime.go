package mal

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

// AnimeEntry holds values such as score, episode and status that we want our
// anime entry to have when we add or update it on our list.
//
// Status is required and can be:
// 1/watching, 2/completed, 3/onhold, 4/dropped, 6/plantowatch
//
// DateStart and DateFinish require 'mmddyyyy' format
//
// EnableDiscussion and EnableRewatching can be: 1=enable, 0=disable
//
// Tags are comma separated: test tag, 2nd tag
type AnimeEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Episode            int      `xml:"episode,omitempty"`
	Status             string   `xml:"status,omitempty"`
	Score              int      `xml:"score,omitempty"`
	DownloadedEpisodes int      `xml:"downloaded_episodes,omitempty"`
	StorageType        int      `xml:"storage_type,omitempty"`
	StorageValue       float64  `xml:"storage_value,omitempty"`
	TimesRewatched     int      `xml:"times_rewatched,omitempty"`
	RewatchValue       int      `xml:"rewatch_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`
	DateFinish         string   `xml:"date_finish,omitempty"`
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"`
	EnableRewatching   int      `xml:"enable_rewatching,omitempty"`
	Comments           string   `xml:"comments,omitempty"`
	FansubGroup        string   `xml:"fansub_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"`
}

type AnimeService struct {
	client *Client
}

func (s *AnimeService) Add(animeID int, entry AnimeEntry) (*http.Response, error) {

	const endpoint = "api/animelist/add/"

	return s.client.post(endpoint, animeID, entry)
}

func (s *AnimeService) Update(animeID int, entry AnimeEntry) (*http.Response, error) {

	const endpoint = "api/animelist/update/"

	return s.client.post(endpoint, animeID, entry)
}

func (s *AnimeService) Delete(animeID int) (*http.Response, error) {

	const endpoint = "api/animelist/delete/"

	return s.client.delete(endpoint, animeID)
}

type AnimeResult struct {
	Rows []AnimeRow `xml:"entry"`
}

type AnimeRow struct {
	Row
	Episodes int `xml:"episodes"`
}

type Row struct {
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

func (s *AnimeService) Search(query string) (*AnimeResult, *http.Response, error) {

	const endpoint = "api/anime/search.xml?q="

	req, err := s.client.NewRequest("GET", fmt.Sprintf("%s%s", endpoint, url.QueryEscape(query)), nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AnimeResult)
	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
