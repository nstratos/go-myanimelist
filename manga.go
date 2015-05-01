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
	client *Client
}

func (s *MangaService) Add(mangaID int, entry MangaEntry) (*Response, error) {

	const endpoint = "api/mangalist/add/"

	return s.client.post(endpoint, mangaID, entry)
}

func (s *MangaService) Update(mangaID int, entry MangaEntry) (*Response, error) {

	const endpoint = "api/mangalist/update/"

	return s.client.post(endpoint, mangaID, entry)
}

func (s *MangaService) Delete(mangaID int) (*Response, error) {

	const endpoint = "api/mangalist/delete/"

	return s.client.delete(endpoint, mangaID)
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

	const endpoint = "api/manga/search.xml?q="
	u := fmt.Sprintf("%s%s", endpoint, url.QueryEscape(query))

	result := new(MangaResult)
	resp, err := s.client.query(u, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}
