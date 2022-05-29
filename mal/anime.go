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
	ID                     int                `json:"id"`
	Title                  string             `json:"title"`
	MainPicture            Picture            `json:"main_picture"`
	AlternativeTitles      Titles             `json:"alternative_titles"`
	StartDate              string             `json:"start_date"`
	EndDate                string             `json:"end_date"`
	Synopsis               string             `json:"synopsis"`
	Mean                   float64            `json:"mean"`
	Rank                   int                `json:"rank"`
	Popularity             int                `json:"popularity"`
	NumListUsers           int                `json:"num_list_users"`
	NumScoringUsers        int                `json:"num_scoring_users"`
	NSFW                   string             `json:"nsfw"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
	MediaType              string             `json:"media_type"`
	Status                 string             `json:"status"`
	Genres                 []Genre            `json:"genres"`
	MyListStatus           AnimeListStatus    `json:"my_list_status"`
	NumEpisodes            int                `json:"num_episodes"`
	StartSeason            StartSeason        `json:"start_season"`
	Broadcast              Broadcast          `json:"broadcast"`
	Source                 string             `json:"source"`
	AverageEpisodeDuration int                `json:"average_episode_duration"`
	Rating                 string             `json:"rating"`
	Pictures               []Picture          `json:"pictures"`
	Background             string             `json:"background"`
	RelatedAnime           []RelatedAnime     `json:"related_anime"`
	RelatedManga           []RelatedManga     `json:"related_manga"`
	Recommendations        []RecommendedAnime `json:"recommendations"`
	Studios                []Studio           `json:"studios"`
	Statistics             Statistics         `json:"statistics"`
}

// Picture is a representative picture from the show.
type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// Titles of the anime in English and Japanese.
type Titles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

// The Genre of the anime.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// The Studio that created the anime.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Status of the user's anime list contained in statistics.
type Status struct {
	Watching    string `json:"watching"`
	Completed   string `json:"completed"`
	OnHold      string `json:"on_hold"`
	Dropped     string `json:"dropped"`
	PlanToWatch string `json:"plan_to_watch"`
}

// Statistics about the anime.
type Statistics struct {
	Status       Status `json:"status"`
	NumListUsers int    `json:"num_list_users"`
}

// RecommendedAnime is a recommended anime.
type RecommendedAnime struct {
	Node               Anime `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// RelatedAnime contains a related anime.
type RelatedAnime struct {
	Node                  Anime  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// StartSeason is the season an anime starts.
type StartSeason struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

// Broadcast is the day and time that the show broadcasts.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
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

// NSFW is an option which sets the NSFW query option. By default this is set to
// false.
type NSFW bool

func (n NSFW) seasonalAnimeApply(v *url.Values) { n.apply(v) }
func (n NSFW) animeListApply(v *url.Values)     { n.apply(v) }
func (n NSFW) mangaListApply(v *url.Values)     { n.apply(v) }
func (n NSFW) apply(v *url.Values)              { v.Set("nsfw", strconv.FormatBool(bool(n))) }

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
