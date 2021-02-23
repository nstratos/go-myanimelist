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

// Limit the results returned by a request.
func Limit(limit int) func(q *url.Values) {
	return func(q *url.Values) {
		q.Set("limit", strconv.Itoa(limit))
	}
}

// Offset request results by an amount.
func Offset(offset int) func(q *url.Values) {
	return func(q *url.Values) {
		q.Set("offset", strconv.Itoa(offset))
	}
}

// Fields allows to choose the fields that should be returned as by default, the
// API doesn't return all fields.
//
// Example:
//
//     Fields("synopsis", "my_list_status{priority,comments}")
func Fields(fields ...string) func(q *url.Values) {
	return func(q *url.Values) {
		if len(fields) != 0 {
			q.Set("fields", strings.Join(fields, ","))
		}
	}
}

// List allows an authenticated user to receive their anime list.
func (s *AnimeService) List(ctx context.Context, search string, options ...func(q *url.Values)) ([]Anime, *Response, error) {
	options = append(options, queryOpt(search))
	return s.client.animeList(ctx, "anime", options...)
}

func queryOpt(query string) func(q *url.Values) {
	return func(q *url.Values) {
		q.Set("q", query)
	}
}

// AnimeRanking allows to choose how the anime will be ranked.
type AnimeRanking string

// Possible AnimeRanking values.
const (
	RankingAll          AnimeRanking = "all"          // Top Anime Series.
	RankingAiring       AnimeRanking = "airing"       // Top Airing Anime.
	RankingUpcoming     AnimeRanking = "upcoming"     // Top Upcoming Anime.
	RankingTV           AnimeRanking = "tv"           // Top Anime TV Series.
	RankingOVA          AnimeRanking = "ova"          // Top Anime OVA Series.
	RankingMovie        AnimeRanking = "movie"        // Top Anime Movies.
	RankingSpecial      AnimeRanking = "special"      // Top Anime Specials.
	RankingByPopularity AnimeRanking = "bypopularity" // Top Anime by Popularity.
	RankingFavorite     AnimeRanking = "favorite"     // Top Favorited Anime.
)

// Ranking allows an authenticated user to receive the top anime based on a
// certain ranking.
func (s *AnimeService) Ranking(ctx context.Context, ranking AnimeRanking, options ...func(q *url.Values)) ([]Anime, *Response, error) {
	options = append(options, rankingOpt(ranking))
	return s.client.animeList(ctx, "anime/ranking", options...)
}

func rankingOpt(ranking AnimeRanking) func(q *url.Values) {
	return func(q *url.Values) {
		q.Set("ranking_type", string(ranking))
	}
}

// AnimeSeason is the airing season of the anime.
type AnimeSeason string

// Possible AnimeSeason values.
const (
	AnimeSeasonWinter AnimeSeason = "winter" // January, February, March.
	AnimeSeasonSpring AnimeSeason = "spring" // April, May, June.
	AnimeSeasonSummer AnimeSeason = "summer" // July, August, September.
	AnimeSeasonFall   AnimeSeason = "fall"   // October, November, December.
)

// SortSeasonalAnimeBy shows the ways the anime results can be sorted.
type SortSeasonalAnimeBy string

// Possible values for sorting seasonal anime.
const (
	ByAnimeScore        SortSeasonalAnimeBy = "anime_score"          // Descending
	ByAnimeNumListUsers SortSeasonalAnimeBy = "anime_num_list_users" // Descending
)

// SeasonalAnimeOption are options specific to the AnimeService.Seasonal method.
type SeasonalAnimeOption func(q *url.Values)

// SortSeasonalAnime allows to choose how the anime results will be sorted when
// using the AnimeService.Seasonal method.
func SortSeasonalAnime(sort SortSeasonalAnimeBy) SeasonalAnimeOption {
	return func(q *url.Values) {
		q.Set("sort", string(sort))
	}
}

// Seasonal allows an authenticated user to receive the seasonal anime by
// providing the year and season.
func (s *AnimeService) Seasonal(ctx context.Context, year int, season AnimeSeason, options ...SeasonalAnimeOption) ([]Anime, *Response, error) {
	oo := make([]func(q *url.Values), len(options))
	for i := range options {
		oo[i] = options[i]
	}
	return s.client.animeList(ctx, fmt.Sprintf("anime/season/%d/%s", year, season), oo...)
}
