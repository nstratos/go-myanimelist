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

// Details returns details about an anime. By default, few anime fields are
// populated. Use the Fields option to specify which fields should be included.
func (s *AnimeService) Details(ctx context.Context, animeID int, options ...DetailsOption) (*Anime, *Response, error) {
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

func (l Limit) pagingApply(v *url.Values)        { l.apply(v) }
func (l Limit) topicsApply(v *url.Values)        { l.apply(v) }
func (l Limit) seasonalAnimeApply(v *url.Values) { l.apply(v) }
func (l Limit) animeListApply(v *url.Values)     { l.apply(v) }
func (l Limit) mangaListApply(v *url.Values)     { l.apply(v) }
func (l Limit) apply(v *url.Values)              { v.Set("limit", strconv.Itoa(int(l))) }

// Offset is an option that sets the offset of the results.
type Offset int

func (o Offset) pagingApply(v *url.Values)        { o.apply(v) }
func (o Offset) topicsApply(v *url.Values)        { o.apply(v) }
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
func (f Fields) myInfoApply(v *url.Values)        { f.apply(v) }
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

// List allows an authenticated user to search and list anime data. You may get
// user specific data by using the optional field "my_list_status".
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

const (
	// AnimeRankingAll returns the top anime series.
	AnimeRankingAll AnimeRanking = "all"
	// AnimeRankingAiring returns the top airing anime.
	AnimeRankingAiring AnimeRanking = "airing"
	// AnimeRankingUpcoming returns the top upcoming anime.
	AnimeRankingUpcoming AnimeRanking = "upcoming"
	// AnimeRankingTV returns the top Anime TV series.
	AnimeRankingTV AnimeRanking = "tv"
	// AnimeRankingOVA returns the top anime OVA series.
	AnimeRankingOVA AnimeRanking = "ova"
	// AnimeRankingMovie returns the top anime movies.
	AnimeRankingMovie AnimeRanking = "movie"
	// AnimeRankingSpecial returns the top anime specials.
	AnimeRankingSpecial AnimeRanking = "special"
	// AnimeRankingByPopularity returns the top anime by popularity.
	AnimeRankingByPopularity AnimeRanking = "bypopularity"
	// AnimeRankingFavorite returns the top favorite Anime.
	AnimeRankingFavorite AnimeRanking = "favorite"
)

func optionFromAnimeRanking(r AnimeRanking) optionFunc {
	return optionFunc(func(v *url.Values) {
		v.Set("ranking_type", string(r))
	})
}

// Ranking allows an authenticated user to receive the top anime based on a
// certain ranking.
func (s *AnimeService) Ranking(ctx context.Context, ranking AnimeRanking, options ...Option) ([]Anime, *Response, error) {
	options = append(options, optionFromAnimeRanking(ranking))
	return s.list(ctx, "anime/ranking", options...)
}

// AnimeSeason is the airing season of the anime.
type AnimeSeason string

const (
	// AnimeSeasonWinter is the winter season of January, February and March.
	AnimeSeasonWinter AnimeSeason = "winter"
	// AnimeSeasonSpring is the spring season of April, May and June.
	AnimeSeasonSpring AnimeSeason = "spring"
	// AnimeSeasonSummer is the summer season of July, August and September.
	AnimeSeasonSummer AnimeSeason = "summer"
	// AnimeSeasonFall is the fall season of October, November and December.
	AnimeSeasonFall AnimeSeason = "fall"
)

// SortSeasonalAnime is an option that allows to sort the anime results.
type SortSeasonalAnime string

const (
	// SortSeasonalByAnimeScore sorts seasonal results by anime score in
	// descending order.
	SortSeasonalByAnimeScore SortSeasonalAnime = "anime_score"
	// SortSeasonalByAnimeNumListUsers sorts seasonal results by anime num list
	// users in descending order.
	SortSeasonalByAnimeNumListUsers SortSeasonalAnime = "anime_num_list_users"
)

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
func (s *AnimeService) Seasonal(ctx context.Context, year int, season AnimeSeason, options ...SeasonalAnimeOption) ([]Anime, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromSeasonalAnimeOption(options[i])
	}
	return s.list(ctx, fmt.Sprintf("anime/season/%d/%s", year, season), oo...)
}

// Suggested returns suggested anime for the authorized user. If the user is new
// comer, this endpoint returns an empty list.
func (s *AnimeService) Suggested(ctx context.Context, options ...Option) ([]Anime, *Response, error) {
	return s.list(ctx, "anime/suggestions", options...)
}
