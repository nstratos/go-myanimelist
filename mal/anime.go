package mal

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

// AnimeEntry represents the values that an anime will have on the list when
// added or updated. Status is required.
type AnimeEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Episode            int      `xml:"episode"`
	Status             string   `xml:"status,omitempty"` // 1|watching, 2|completed, 3|onhold, 4|dropped, 6|plantowatch
	Score              int      `xml:"score"`
	DownloadedEpisodes int      `xml:"downloaded_episodes,omitempty"`
	StorageType        int      `xml:"storage_type,omitempty"`
	StorageValue       float64  `xml:"storage_value,omitempty"`
	TimesRewatched     int      `xml:"times_rewatched"`
	RewatchValue       int      `xml:"rewatch_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`  // mmddyyyy
	DateFinish         string   `xml:"date_finish,omitempty"` // mmddyyyy
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"` // 1=enable, 0=disable
	EnableRewatching   int      `xml:"enable_rewatching"`           // 1=enable, 0=disable
	Comments           string   `xml:"comments"`
	FansubGroup        string   `xml:"fansub_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"` // comma separated: test tag, 2nd tag
}

// AnimeService handles communication with the Anime List methods of the
// MyAnimeList API.
//
// MyAnimeList API docs: http://myanimelist.net/modules.php?go=api
type AnimeService struct {
	client         *Client
	AddEndpoint    *url.URL
	UpdateEndpoint *url.URL
	DeleteEndpoint *url.URL
	SearchEndpoint *url.URL
	ListEndpoint   *url.URL
}

// Add allows an authenticated user to add an anime to their anime list.
func (s *AnimeService) Add(animeID int, entry AnimeEntry) (*Response, error) {

	return s.client.post(s.AddEndpoint.String(), animeID, entry)
}

// Update allows an authenticated user to update an anime on their anime list.
//
// Note: MyAnimeList.net updates the MyLastUpdated value of an Anime only if it
// receives an Episode update change. For example:
//
//    Updating Episode 0 -> 1 will update MyLastUpdated
//    Updating Episode 0 -> 0 will not update MyLastUpdated
//    Updating Status  1 -> 2 will not update MyLastUpdated
//    Updating Rating  5 -> 8 will not update MyLastUpdated
//
// As a consequence, you might perform a number of updates on a certain anime
// that will not affect it's MyLastUpdate unless one of the updates happens to
// change the episode.  This behavior is important to know if your application
// performs updates and cares about when an anime was last updated.
func (s *AnimeService) Update(animeID int, entry AnimeEntry) (*Response, error) {

	return s.client.post(s.UpdateEndpoint.String(), animeID, entry)
}

// Delete allows an authenticated user to delete an anime from their anime list.
func (s *AnimeService) Delete(animeID int) (*Response, error) {

	return s.client.delete(s.DeleteEndpoint.String(), animeID)
}

// AnimeResult represents the result of an anime search.
type AnimeResult struct {
	Rows []AnimeRow `xml:"entry"`
}

// AnimeRow represents each row of an anime search result.
type AnimeRow struct {
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
	Episodes  int     `xml:"episodes"`
}

// Search allows an authenticated user to search anime titles. Upon failure it
// will return ErrNoContent.
func (s *AnimeService) Search(query string) (*AnimeResult, *Response, error) {

	v := s.SearchEndpoint.Query()
	v.Set("q", query)
	s.SearchEndpoint.RawQuery = v.Encode()

	result := new(AnimeResult)
	resp, err := s.client.get(s.SearchEndpoint.String(), result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// AnimeList represents the anime list of a user.
type AnimeList struct {
	MyInfo AnimeMyInfo `xml:"myinfo"`
	Anime  []Anime     `xml:"anime"`
	Error  string      `xml:"error"`
}

// AnimeMyInfo represents the user's info (like number of watching, completed etc)
// that is returned when requesting his/her anime list.
type AnimeMyInfo struct {
	ID                int    `xml:"user_id"`
	Name              string `xml:"user_name"`
	Completed         int    `xml:"user_completed"`
	OnHold            int    `xml:"user_onhold"`
	Dropped           int    `xml:"user_dropped"`
	DaysSpentWatching string `xml:"user_days_spent_watching"`
	Watching          int    `xml:"user_watching"`
	PlanToWatch       int    `xml:"user_plantowatch"`
}

// Anime represents an anime from MyAnimeList with it's data contained in the
// fields starting with Series. It also contains user specific fields
// that start with My (for example the score the user has set for that anime).
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
	MyStatus          int    `xml:"my_status"`     // 1 = watching, 2 = completed, 3 = onhold, 4 = dropped, 6 = plantowatch
	MyRewatching      int    `xml:"my_rewatching"` // Officially undocumented but it seems that 1=true and 0=false.
	MyRewatchingEp    int    `xml:"my_rewatching_ep"`
	MyLastUpdated     string `xml:"my_last_updated"`
	MyTags            string `xml:"my_tags"`
	MyWatchedEpisodes int    `xml:"my_watched_episodes"`
}

// List allows an authenticated user to receive the anime list of a user.
func (s *AnimeService) List(username string) (*AnimeList, *Response, error) {

	v := s.ListEndpoint.Query()
	v.Set("status", "all")
	v.Set("type", "anime")
	v.Set("u", username)
	s.ListEndpoint.RawQuery = v.Encode()

	list := new(AnimeList)
	resp, err := s.client.get(s.ListEndpoint.String(), list)
	if err != nil {
		return nil, resp, err
	}

	if list.Error != "" {
		return list, resp, fmt.Errorf("%v", list.Error)
	}

	return list, resp, nil
}
