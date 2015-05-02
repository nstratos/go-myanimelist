package mal

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

// MangaEntry holds values such as score, chapter and status that we want our
// manga entry to have when we add or update it on our list.
//
// Status is required and can be:
// 1/reading, 2/completed, 3/onhold, 4/dropped, 6/plantoread
//
// DateStart and DateFinish require 'mmddyyyy' format
//
// EnableDiscussion and EnableRereading can be: 1=enable, 0=disable
//
// Tags are comma separated: test tag, 2nd tag
type MangaEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Volume             int      `xml:"volume,omitempty"`
	Chapter            int      `xml:"chapter,omitempty"`
	Status             string   `xml:"status,omitempty"`
	Score              int      `xml:"score,omitempty"`
	DownloadedChapters int      `xml:"downloaded_chapters,omitempty"`
	TimesReread        int      `xml:"times_reread,omitempty"`
	RereadValue        int      `xml:"reread_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`
	DateFinish         string   `xml:"date_finish,omitempty"`
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"`
	EnableRereading    int      `xml:"enable_rereading,omitempty"`
	Comments           string   `xml:"comments,omitempty"`
	ScanGroup          string   `xml:"scan_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"`
	RetailVolumes      int      `xml:"retail_volumes,omitempty"`
}

type MangaService struct {
	client         *Client
	AddEndpoint    *url.URL
	UpdateEndpoint *url.URL
	DeleteEndpoint *url.URL
	SearchEndpoint *url.URL
	ListEndpoint   *url.URL
}

func (s *MangaService) Add(mangaID int, entry MangaEntry) (*Response, error) {

	return s.client.post(s.AddEndpoint.String(), mangaID, entry)
}

func (s *MangaService) Update(mangaID int, entry MangaEntry) (*Response, error) {

	return s.client.post(s.UpdateEndpoint.String(), mangaID, entry)
}

func (s *MangaService) Delete(mangaID int) (*Response, error) {

	return s.client.delete(s.DeleteEndpoint.String(), mangaID)
}

type MangaResult struct {
	Rows []MangaRow `xml:"entry"`
}

type MangaRow struct {
	Row
	Chapters int `xml:"chapters"`
	Volumes  int `xml:"volumes"`
}

func (s *MangaService) Search(query string) (*MangaResult, *Response, error) {

	u := fmt.Sprintf("%s?q=%s", s.SearchEndpoint, url.QueryEscape(query))

	result := new(MangaResult)
	resp, err := s.client.query(u, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

type MangaList struct {
	MyInfo MyMangaInfo `xml:"myinfo"`
	Manga  []Manga     `xml:"manga"`
	Error  string      `xml:"error"`
}

type MyMangaInfo struct {
	ID                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
	Reading           int    `xml:"user_reading"`
	PlanToRead        int    `xml:"user_plantoread"`
}

// Manga holds data for each manga entry. User specific data for each manga
// are also held in the fields starting with My.
//
// MyStatus: 1 = watching, 2 = completed, 3 = on hold, 4 = dropped, 6 = plantowatch
type Manga struct {
	SeriesMangaDBID int    `xml:"series_mangadb_id"`
	SeriesChapters  int    `xml:"series_chapters"`
	SeriesVolumes   int    `xml:"series_volumes"`
	SeriesTitle     string `xml:"series_title"`
	SeriesSynonyms  string `xml:"series_synonyms"`
	SeriesType      int    `xml:"series_type"`
	SeriesStatus    int    `xml:"series_status"`
	SeriesStart     string `xml:"series_start"`
	SeriesEnd       string `xml:"series_end"`
	SeriesImage     string `xml:"series_image"`
	MyID            int    `xml:"my_id"`
	MyStartDate     string `xml:"my_start_date"`
	MyFinishDate    string `xml:"my_finish_date"`
	MyScore         int    `xml:"my_score"`
	MyStatus        int    `xml:"my_status"`
	MyRewatching    string `xml:"my_rewatching"`
	MyRewatchingEp  int    `xml:"my_rewatching_ep"`
	MyLastUpdated   string `xml:"my_last_updated"`
	MyTags          string `xml:"my_tags"`
	MyReadChapters  int    `xml:"my_read_chapters"`
	MyReadVolumes   int    `xml:"my_read_volumes"`
}

func (s *MangaService) List(username string) (*MangaList, *Response, error) {

	u := fmt.Sprintf("%s?status=all&type=manga&u=%s", s.ListEndpoint, url.QueryEscape(username))

	list := new(MangaList)
	resp, err := s.client.query(u, list)
	if err != nil {
		return nil, resp, err
	}

	if list.Error != "" {
		return list, resp, fmt.Errorf("%v", list.Error)
	}

	return list, resp, nil
}
