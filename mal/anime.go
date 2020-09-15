package mal

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// AnimeEntry represents the values that an anime will have on the list when
// added or updated. Status is required.
type AnimeEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Episode            int      `xml:"episode"`
	Status             Status   `xml:"status,omitempty"` // Use the package constants: mal.Current, mal.Completed, etc.
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
	Tags               string   `xml:"tags,omitempty"` // comma separated
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

// Anime represents a MyAnimeList anime.
type Anime struct {
	ID                     int              `json:"id"`
	Title                  string           `json:"title"`
	MainPicture            Picture          `json:"main_picture"`
	AlternativeTitles      Titles           `json:"alternative_titles"`
	StartDate              string           `json:"start_date"`
	EndDate                string           `json:"end_date"`
	Synopsis               string           `json:"synopsis"`
	Mean                   float64          `json:"mean"`
	Rank                   int              `json:"rank"`
	Popularity             int              `json:"popularity"`
	NumListUsers           int              `json:"num_list_users"`
	NumScoringUsers        int              `json:"num_scoring_users"`
	NSFW                   string           `json:"nsfw"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
	MediaType              string           `json:"media_type"`
	Status                 string           `json:"status"`
	Genres                 []Genre          `json:"genres"`
	MyListStatus           MyListStatus     `json:"my_list_status"`
	NumEpisodes            int              `json:"num_episodes"`
	StartSeason            StartSeason      `json:"start_season"`
	Broadcast              Broadcast        `json:"broadcast"`
	Source                 string           `json:"source"`
	AverageEpisodeDuration int              `json:"average_episode_duration"`
	Rating                 string           `json:"rating"`
	Pictures               []Picture        `json:"pictures"`
	Background             string           `json:"background"`
	RelatedAnime           []RelatedAnime   `json:"related_anime"`
	RelatedManga           []interface{}    `json:"related_manga"`
	Recommendations        []Recommendation `json:"recommendations"`
	Studios                []Studio         `json:"studios"`
	Statistics             Statistics       `json:"statistics"`
}

type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type Titles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MyListStatus struct {
	Status             string    `json:"status"`
	Score              int       `json:"score"`
	NumEpisodesWatched int       `json:"num_episodes_watched"`
	IsRewatching       bool      `json:"is_rewatching"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	Watching    string `json:"watching"`
	Completed   string `json:"completed"`
	OnHold      string `json:"on_hold"`
	Dropped     string `json:"dropped"`
	PlanToWatch string `json:"plan_to_watch"`
}

type Statistics struct {
	Status       Status `json:"status"`
	NumListUsers int    `json:"num_list_users"`
}

type Recommendation struct {
	Node               Node `json:"node"`
	NumRecommendations int  `json:"num_recommendations"`
}

type Node struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	MainPicture Picture `json:"main_picture"`
}

type RelatedAnime struct {
	Node                  Node   `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

type StartSeason struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
}

// Details returns details about an anime.
func (s *AnimeService) Details(ctx context.Context, id int64) (*Anime, *Response, error) {
	var u string

	u = fmt.Sprintf("anime/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Anime)
	resp, err := s.client.Do(ctx, req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, nil
}

// AnimeList represents the anime list of a user.
type animeList struct {
	Data []struct {
		Anime Anime `json:"node"`
	}
	Paging Paging `json:"paging"`
}

type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// List allows an authenticated user to receive the anime list of a user.
func (s *AnimeService) List(ctx context.Context, query string, limit, offset int, fields ...string) ([]Anime, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "anime", nil)
	if err != nil {
		return nil, nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset != 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	if len(fields) != 0 {
		q.Set("fields", strings.Join(fields, ","))

	}
	req.URL.RawQuery = q.Encode()

	list := new(animeList)
	resp, err := s.client.Do(ctx, req, list)
	if err != nil {
		return nil, resp, err
	}

	if list.Paging.Previous != "" {
		offset, err := parseOffset(list.Paging.Previous)
		if err != nil {
			return nil, resp, fmt.Errorf("previous: %s", err)
		}
		resp.PrevOffset = offset
	}
	if list.Paging.Next != "" {
		offset, err := parseOffset(list.Paging.Next)
		if err != nil {
			return nil, resp, fmt.Errorf("next: %s", err)
		}
		resp.NextOffset = offset
	}
	anime := []Anime{}
	for _, d := range list.Data {
		anime = append(anime, d.Anime)
	}

	return anime, resp, nil
}

func parseOffset(urlStr string) (int, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return 0, fmt.Errorf("parsing URL %q: %s", urlStr, err)
	}
	offset, err := strconv.Atoi(u.Query().Get("offset"))
	if err != nil {
		return 0, fmt.Errorf("parsing offset: %s", err)
	}
	return offset, nil
}

// Add allows an authenticated user to add an anime to their anime list.
func (s *AnimeService) Add(animeID int, entry AnimeEntry) (*Response, error) {

	return s.client.post(s.AddEndpoint.String(), animeID, entry, true)
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
// that will not affect its MyLastUpdate unless one of the updates happens to
// change the episode. This behavior is important to know if your application
// performs updates and cares about when an anime was last updated.
func (s *AnimeService) Update(animeID int, entry AnimeEntry) (*Response, error) {

	return s.client.post(s.UpdateEndpoint.String(), animeID, entry, true)
}

// Delete allows an authenticated user to delete an anime from their anime list.
func (s *AnimeService) Delete(animeID int) (*Response, error) {

	return s.client.delete(s.DeleteEndpoint.String(), animeID, true)
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

// Search allows an authenticated user to search anime titles. If nothing is
// found, it will return an ErrNoContent error.
func (s *AnimeService) Search(query string) (*AnimeResult, *Response, error) {

	v := s.SearchEndpoint.Query()
	v.Set("q", query)
	s.SearchEndpoint.RawQuery = v.Encode()

	result := new(AnimeResult)
	resp, err := s.client.get(s.SearchEndpoint.String(), result, true)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// AnimeMyInfo represents the user's info which contains stats about the anime
// that exist in their anime list. For example how many anime they have
// completed, how many anime they are currently watching etc. It is returned as
// part of their AnimeList.
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

// Anime2 represents a MyAnimeList anime. The data of the anime are stored in
// the fields starting with Series. User specific data are stored in the fields
// starting with My. For example, the score the user has set for that anime is
// stored in the MyScore field.
type Anime2 struct {
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
	MyStatus          Status `xml:"my_status"` // Use the package constants: mal.Current, mal.Completed, etc.
	MyRewatching      int    `xml:"my_rewatching"`
	MyRewatchingEp    int    `xml:"my_rewatching_ep"`
	MyLastUpdated     string `xml:"my_last_updated"`
	MyTags            string `xml:"my_tags"`
	MyWatchedEpisodes int    `xml:"my_watched_episodes"`
}
