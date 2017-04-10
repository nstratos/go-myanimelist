package mal

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

// MangaEntry represents the values that an manga will have on the list when
// added or updated. Status is required.
type MangaEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Volume             int      `xml:"volume,omitempty"`
	Chapter            int      `xml:"chapter,omitempty"`
	Status             int      `xml:"status,omitempty"` // Use the package constants: StatusReading, StatusCompleted, etc.
	Score              int      `xml:"score,omitempty"`
	DownloadedChapters int      `xml:"downloaded_chapters,omitempty"`
	TimesReread        int      `xml:"times_reread,omitempty"`
	RereadValue        int      `xml:"reread_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`  // mmddyyyy
	DateFinish         string   `xml:"date_finish,omitempty"` // mmddyyyy
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"` // 1=enable, 0=disable
	EnableRereading    int      `xml:"enable_rereading,omitempty"`  // 1=enable, 0=disable
	Comments           string   `xml:"comments,omitempty"`
	ScanGroup          string   `xml:"scan_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"` // comma separated: test tag, 2nd tag
	RetailVolumes      int      `xml:"retail_volumes,omitempty"`
}

// MangaService handles communication with the Manga List methods of the
// MyAnimeList API.
//
// MyAnimeList API docs: http://myanimelist.net/modules.php?go=api
type MangaService struct {
	client         *Client
	AddEndpoint    *url.URL
	UpdateEndpoint *url.URL
	DeleteEndpoint *url.URL
	SearchEndpoint *url.URL
	ListEndpoint   *url.URL
}

// Add allows an authenticated user to add a manga to their manga list.
func (s *MangaService) Add(mangaID int, entry MangaEntry) (*Response, error) {

	return s.client.post(s.AddEndpoint.String(), mangaID, entry)
}

// Update allows an authenticated user to update an manga on their manga list.
func (s *MangaService) Update(mangaID int, entry MangaEntry) (*Response, error) {

	return s.client.post(s.UpdateEndpoint.String(), mangaID, entry)
}

// Delete allows an authenticated user to delete an manga from their manga list.
func (s *MangaService) Delete(mangaID int) (*Response, error) {

	return s.client.delete(s.DeleteEndpoint.String(), mangaID)
}

// MangaResult represents the result of an manga search.
type MangaResult struct {
	Rows []MangaRow `xml:"entry"`
}

// MangaRow represents each row of an manga search result.
type MangaRow struct {
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
	Chapters  int     `xml:"chapters"`
	Volumes   int     `xml:"volumes"`
}

// Search allows an authenticated user to search manga titles. If nothing is
// found, it will return an ErrNoContent error.
func (s *MangaService) Search(query string) (*MangaResult, *Response, error) {

	v := s.SearchEndpoint.Query()
	v.Set("q", query)
	s.SearchEndpoint.RawQuery = v.Encode()

	result := new(MangaResult)
	resp, err := s.client.get(s.SearchEndpoint.String(), result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// MangaList represents the manga list of a user.
type MangaList struct {
	MyInfo MangaMyInfo `xml:"myinfo"`
	Manga  []Manga     `xml:"manga"`
	Error  string      `xml:"error"`
}

// MangaMyInfo represents the user's info which contains stats about the manga
// that exist in their manga list. For example how many manga they have
// completed, how many manga they are currently reading etc. It is returned as
// part of their MangaList.
type MangaMyInfo struct {
	ID                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
	Reading           int    `xml:"user_reading"`
	PlanToRead        int    `xml:"user_plantoread"`
}

// Manga represents a MyAnimeList manga. The data of the manga are stored in
// the fields starting with Series. User specific data are stored in the fields
// starting with My. For example, the score the user has set for that manga is
// stored in the MyScore field.
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
	MyStatus        int    `xml:"my_status"`     // 1 = reading, 2 = completed, 3 = onhold, 4 = dropped, 6 = plantoread
	MyRereading     int    `xml:"my_rereadingg"` // MyAnimeList spells it my_rereadingg. Possible values seem to be 1=true and 0=false.
	MyRereadingChap int    `xml:"my_rereading_chap"`
	MyLastUpdated   string `xml:"my_last_updated"`
	MyTags          string `xml:"my_tags"`
	MyReadChapters  int    `xml:"my_read_chapters"`
	MyReadVolumes   int    `xml:"my_read_volumes"`
}

// List allows an authenticated user to receive the manga list of a user.
func (s *MangaService) List(username string) (*MangaList, *Response, error) {

	v := s.ListEndpoint.Query()
	v.Set("status", "all")
	v.Set("type", "manga")
	v.Set("u", username)
	s.ListEndpoint.RawQuery = v.Encode()

	list := new(MangaList)
	resp, err := s.client.get(s.ListEndpoint.String(), list)
	if err != nil {
		return nil, resp, err
	}

	if list.Error != "" {
		return list, resp, fmt.Errorf("%v", list.Error)
	}

	return list, resp, nil
}
