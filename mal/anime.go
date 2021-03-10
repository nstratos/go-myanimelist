package mal

import (
	"context"
	"fmt"
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
	ID                     int                `json:"id,omitempty"`
	Title                  string             `json:"title,omitempty"`
	MainPicture            Picture            `json:"main_picture,omitempty"`
	AlternativeTitles      Titles             `json:"alternative_titles,omitempty"`
	StartDate              string             `json:"start_date,omitempty"`
	EndDate                string             `json:"end_date,omitempty"`
	Synopsis               string             `json:"synopsis,omitempty"`
	Mean                   float64            `json:"mean,omitempty"`
	Rank                   int                `json:"rank,omitempty"`
	Popularity             int                `json:"popularity,omitempty"`
	NumListUsers           int                `json:"num_list_users,omitempty"`
	NumScoringUsers        int                `json:"num_scoring_users,omitempty"`
	NSFW                   string             `json:"nsfw,omitempty"`
	CreatedAt              time.Time          `json:"created_at,omitempty"`
	UpdatedAt              time.Time          `json:"updated_at,omitempty"`
	MediaType              string             `json:"media_type,omitempty"`
	Status                 string             `json:"status,omitempty"`
	Genres                 []Genre            `json:"genres,omitempty"`
	MyListStatus           AnimeListStatus    `json:"my_list_status,omitempty"`
	NumEpisodes            int                `json:"num_episodes,omitempty"`
	StartSeason            StartSeason        `json:"start_season,omitempty"`
	Broadcast              Broadcast          `json:"broadcast,omitempty"`
	Source                 string             `json:"source,omitempty"`
	AverageEpisodeDuration int                `json:"average_episode_duration,omitempty"`
	Rating                 string             `json:"rating,omitempty"`
	Pictures               []Picture          `json:"pictures,omitempty"`
	Background             string             `json:"background,omitempty"`
	RelatedAnime           []RelatedAnime     `json:"related_anime,omitempty"`
	RelatedManga           []RelatedManga     `json:"related_manga,omitempty"`
	Recommendations        []RecommendedAnime `json:"recommendations,omitempty"`
	Studios                []Studio           `json:"studios,omitempty"`
	Statistics             Statistics         `json:"statistics,omitempty"`
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

// The Studio that created the anime.
type Studio struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Status of the user's anime list contained in statistics.
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

// RecommendedAnime is a recommended anime.
type RecommendedAnime struct {
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

// DetailsOption is an option specific for the anime and manga details methods.
type DetailsOption interface {
	detailsApply(v *url.Values)
}

// Details returns details about an anime.
func (s *AnimeService) Details(ctx context.Context, animeID int64, options ...DetailsOption) (*Anime, *Response, error) {
	a := new(Anime)
	resp, err := s.client.details(ctx, fmt.Sprintf("anime/%d", animeID), a, options...)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, nil
}

// Option is implemented by types that can be used as options in most methods
// such as Limit, Offset and Fields.
type Option interface {
	apply(v *url.Values)
}

type optionFunc func(v *url.Values)

func (f optionFunc) apply(v *url.Values) {
	f(v)
}

// Limit is an option that limits the results returned by a request.
type Limit int

func (l Limit) seasonalAnimeApply(v *url.Values) { l.apply(v) }
func (l Limit) animeListApply(v *url.Values)     { l.apply(v) }
func (l Limit) mangaListApply(v *url.Values)     { l.apply(v) }
func (l Limit) apply(v *url.Values)              { v.Set("limit", strconv.Itoa(int(l))) }

// Offset is an option that sets the offset of the results.
type Offset int

func (o Offset) seasonalAnimeApply(v *url.Values) { o.apply(v) }
func (o Offset) animeListApply(v *url.Values)     { o.apply(v) }
func (o Offset) mangaListApply(v *url.Values)     { o.apply(v) }
func (o Offset) apply(v *url.Values)              { v.Set("offset", strconv.Itoa(int(o))) }

// Fields is an option that allows to choose the fields that should be returned
// as by default, the API doesn't return all fields.
//
// Example:
//
//     Fields{"synopsis", "my_list_status{priority,comments}"}
type Fields []string

func (f Fields) seasonalAnimeApply(v *url.Values) { f.apply(v) }
func (f Fields) animeListApply(v *url.Values)     { f.apply(v) }
func (f Fields) mangaListApply(v *url.Values)     { f.apply(v) }
func (f Fields) detailsApply(v *url.Values)       { f.apply(v) }
func (f Fields) apply(v *url.Values) {
	if len(f) != 0 {
		v.Set("fields", strings.Join(f, ","))
	}
}

func optionFromQuery(query string) optionFunc {
	return optionFunc(func(v *url.Values) {
		v.Set("q", query)
	})
}

// List allows an authenticated user to receive their anime list.
func (s *AnimeService) List(ctx context.Context, search string, options ...Option) ([]Anime, *Response, error) {
	options = append(options, optionFromQuery(search))
	return s.list(ctx, "anime", options...)
}

func (s *AnimeService) list(ctx context.Context, path string, options ...Option) ([]Anime, *Response, error) {
	list := new(animeList)
	resp, err := s.client.list(ctx, path, list, options...)
	if err != nil {
		return nil, resp, err
	}
	anime := make([]Anime, len(list.Data))
	for i := range list.Data {
		anime[i] = list.Data[i].Anime
	}
	return anime, resp, nil
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

func optionFromAnimeRanking(r AnimeRanking) optionFunc {
	return optionFunc(func(v *url.Values) {
		v.Set("ranking_type", string(r))
	})
}

// type animeRankingOption AnimeRanking

// func (o animeRankingOption) seasonalAnimeApply(v *url.Values) { o.apply(v) }
// func (o animeRankingOption) apply(v *url.Values)              { v.Set("ranking_type", string(o)) }

// Ranking allows an authenticated user to receive the top anime based on a
// certain ranking.
func (s *AnimeService) Ranking(ctx context.Context, ranking AnimeRanking, options ...Option) ([]Anime, *Response, error) {
	options = append(options, optionFromAnimeRanking(ranking))
	return s.list(ctx, "anime/ranking", options...)
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

// SortSeasonalAnime is an option that allows to sort the anime results.
type SortSeasonalAnime string

// Possible values for sorting seasonal anime.
const (
	SortSeasonalByAnimeScore        SortSeasonalAnime = "anime_score"          // Descending
	SortSeasonalByAnimeNumListUsers SortSeasonalAnime = "anime_num_list_users" // Descending
)

//func (s SortSeasonalAnime) apply(v *url.Values)              { s.seasonalAnimeApply(v) }
func (s SortSeasonalAnime) seasonalAnimeApply(v *url.Values) { v.Set("sort", string(s)) }

// SeasonalAnimeOption are options specific to the AnimeService.Seasonal method.
type SeasonalAnimeOption interface {
	seasonalAnimeApply(v *url.Values)
}

func optionFromSeasonalAnimeOption(o SeasonalAnimeOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.seasonalAnimeApply(v)
	})
}

// Seasonal allows an authenticated user to receive the seasonal anime by
// providing the year and season. The results can be sorted using an option.
//
// Example:
//
//     anime, _, err := c.Anime.Seasonal(ctx, 2020, mal.AnimeSeasonFall,
//         mal.Fields{"rank", "popularity"},
//         mal.SortSeasonalByAnimeNumListUsers,
//         mal.Limit(10),
//         mal.Offset(0),
//     )
//     if err != nil {
//         return err
//     }
//     for _, a := range anime {
//         fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
//     }
func (s *AnimeService) Seasonal(ctx context.Context, year int, season AnimeSeason, options ...SeasonalAnimeOption) ([]Anime, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromSeasonalAnimeOption(options[i])
	}
	return s.list(ctx, fmt.Sprintf("anime/season/%d/%s", year, season), oo...)
}

// Suggested returns suggested anime for the authorized user. If the user is new
// comer, this endpoint returns an empty list.
//
// Example:
//
//     anime, _, err := c.Anime.Suggested(ctx,
//         mal.Limit(10),
//         mal.Fields{"rank", "popularity"},
//     )
//     if err != nil {
//         return err
//     }
//     for _, a := range anime {
//         fmt.Printf("Rank: %5d, Popularity: %5d %s\n", a.Rank, a.Popularity, a.Title)
//     }
func (s *AnimeService) Suggested(ctx context.Context, options ...Option) ([]Anime, *Response, error) {
	return s.list(ctx, "anime/suggestions", options...)
}
