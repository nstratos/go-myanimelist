package mal

import (
	"encoding/xml"
	"fmt"
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

func (s *AnimeService) Add(animeID int, entry AnimeEntry) (*Response, error) {

	const endpoint = "api/animelist/add/"

	return s.client.post(endpoint, animeID, entry)
}

func (s *AnimeService) Update(animeID int, entry AnimeEntry) (*Response, error) {

	const endpoint = "api/animelist/update/"

	return s.client.post(endpoint, animeID, entry)
}

func (s *AnimeService) Delete(animeID int) (*Response, error) {

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

func (s *AnimeService) Search(query string) (*AnimeResult, *Response, error) {

	const endpoint = "api/anime/search.xml?q="
	u := fmt.Sprintf("%s%s", endpoint, url.QueryEscape(query))

	result := new(AnimeResult)
	resp, err := s.client.query(u, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

type AnimeList struct {
	MyInfo MyAnimeInfo `xml:"myinfo"`
	Anime  []Anime     `xml:"anime"`
}

type MyAnimeInfo struct {
	ID                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
	Watching          int    `xml:"user_watching"`
	PlanToWatch       int    `xml:"user_plantowatch"`
}

type Anime struct {
	SeriesAnimeDBID   int    `xml:"series_animedb_id"`
	SeriesEpisodes    int    `xml:"series_episodes"`
	SeriesTitle       string `xml:"series_title"`
	SeriesSynonyms    string `xml:"series_synonyms"`
	SeriesType        int    `xml:"series_type"`
	SeriesStatus      int    `xml:"series_status"`
	SeriesStart       string `xml:"series_start"`
	SeriesEnd         string `xml:"series_end"`
	SeriesImage       string `xml:"series_image"`
	MyID              int    `xml:"my_id"`
	MyStartDate       string `xml:"my_start_date"`
	MyFinishDate      string `xml:"my_finish_date"`
	MyScore           int    `xml:"my_score"`
	MyStatus          int    `xml:"my_status"`
	MyRewatching      string `xml:"my_rewatching"`
	MyRewatchingEp    int    `xml:"my_rewatching_ep"`
	MyLastUpdated     string `xml:"my_last_updated"`
	MyTags            string `xml:"my_tags"`
	MyWatchedEpisodes int    `xml:"my_watched_episodes"`
}

func (s *AnimeService) List(username string) (*AnimeList, *Response, error) {

	const endpoint = "malappinfo.php?status=all&type=anime&u="
	u := fmt.Sprintf("%s%s", endpoint, url.QueryEscape(username))

	list := new(AnimeList)
	resp, err := s.client.query(u, list)
	if err != nil {
		return nil, resp, err
	}
	return list, resp, nil
}
