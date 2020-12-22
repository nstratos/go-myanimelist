package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// AnimeService handles communication with the anime related methods of the
// MyAnimeList API:
//
// https://myanimelist.net/apiconfig/references/api/v2#tag/anime
// https://myanimelist.net/apiconfig/references/api/v2#tag/user-animelist
type AnimeService struct {
	client *Client
}

// Anime represents a MyAnimeList anime.
type Anime struct {
	ID                     int              `json:"id,omitempty"`
	Title                  string           `json:"title,omitempty"`
	MainPicture            Picture          `json:"main_picture,omitempty"`
	AlternativeTitles      Titles           `json:"alternative_titles,omitempty"`
	StartDate              string           `json:"start_date,omitempty"`
	EndDate                string           `json:"end_date,omitempty"`
	Synopsis               string           `json:"synopsis,omitempty"`
	Mean                   float64          `json:"mean,omitempty"`
	Rank                   int              `json:"rank,omitempty"`
	Popularity             int              `json:"popularity,omitempty"`
	NumListUsers           int              `json:"num_list_users,omitempty"`
	NumScoringUsers        int              `json:"num_scoring_users,omitempty"`
	NSFW                   string           `json:"nsfw,omitempty"`
	CreatedAt              time.Time        `json:"created_at,omitempty"`
	UpdatedAt              time.Time        `json:"updated_at,omitempty"`
	MediaType              string           `json:"media_type,omitempty"`
	Status                 string           `json:"status,omitempty"`
	Genres                 []Genre          `json:"genres,omitempty"`
	MyListStatus           MyListStatus     `json:"my_list_status,omitempty"`
	NumEpisodes            int              `json:"num_episodes,omitempty"`
	StartSeason            StartSeason      `json:"start_season,omitempty"`
	Broadcast              Broadcast        `json:"broadcast,omitempty"`
	Source                 string           `json:"source,omitempty"`
	AverageEpisodeDuration int              `json:"average_episode_duration,omitempty"`
	Rating                 string           `json:"rating,omitempty"`
	Pictures               []Picture        `json:"pictures,omitempty"`
	Background             string           `json:"background,omitempty"`
	RelatedAnime           []RelatedAnime   `json:"related_anime,omitempty"`
	RelatedManga           []interface{}    `json:"related_manga,omitempty"`
	Recommendations        []Recommendation `json:"recommendations,omitempty"`
	Studios                []Studio         `json:"studios,omitempty"`
	Statistics             Statistics       `json:"statistics,omitempty"`
}

// Picture is a representative picture from the show.
type Picture struct {
	Medium string `json:"medium,omitempty"`
	Large  string `json:"large,omitempty"`
}

// Titles of the anime in English and Japanese.
type Titles struct {
	Synonyms []string `json:"synonyms,omitempty"`
	En       string   `json:"en,omitempty"`
	Ja       string   `json:"ja,omitempty"`
}

// The Genre of the anime.
type Genre struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// MyListStatus is the user's list status.
type MyListStatus struct {
	Status             string    `json:"status,omitempty"`
	Score              int       `json:"score,omitempty"`
	NumEpisodesWatched int       `json:"num_episodes_watched,omitempty"`
	IsRewatching       bool      `json:"is_rewatching,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// The Studio that created the anime.
type Studio struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Status of the anime.
type Status struct {
	Watching    string `json:"watching,omitempty"`
	Completed   string `json:"completed,omitempty"`
	OnHold      string `json:"on_hold,omitempty"`
	Dropped     string `json:"dropped,omitempty"`
	PlanToWatch string `json:"plan_to_watch,omitempty"`
}

// Statistics about the anime.
type Statistics struct {
	Status       Status `json:"status,omitempty"`
	NumListUsers int    `json:"num_list_users,omitempty"`
}

// Recommendation is a recommended anime.
type Recommendation struct {
	Node               Anime `json:"node,omitempty"`
	NumRecommendations int   `json:"num_recommendations,omitempty"`
}

// RelatedAnime contains a related anime.
type RelatedAnime struct {
	Node                  Anime  `json:"node,omitempty"`
	RelationType          string `json:"relation_type,omitempty"`
	RelationTypeFormatted string `json:"relation_type_formatted,omitempty"`
}

// StartSeason is the season an anime starts.
type StartSeason struct {
	Year   int    `json:"year,omitempty"`
	Season string `json:"season,omitempty"`
}

// Broadcast is the day and time that the show broadcasts.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week,omitempty"`
	StartTime    string `json:"start_time,omitempty"`
}

// Details returns details about an anime.
func (s *AnimeService) Details(ctx context.Context, id int64) (*Anime, *Response, error) {
	u := fmt.Sprintf("anime/%d", id)
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

// animeList represents the anime list of a user.
type animeList struct {
	Data []struct {
		Anime Anime `json:"node"`
	}
	Paging Paging `json:"paging"`
}

// Paging provides access to the next and previous page URLs when there are
// pages of results.
type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// List allows an authenticated user to receive their anime list.
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
