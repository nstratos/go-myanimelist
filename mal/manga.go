package mal

import (
	"context"
	"fmt"
	"time"
)

// MangaService handles communication with the manga related methods of the
// MyAnimeList API:
//
// https://myanimelist.net/apiconfig/references/api/v2#tag/manga
// https://myanimelist.net/apiconfig/references/api/v2#tag/user-mangalist
type MangaService struct {
	client *Client
}

// Manga represents a MyAnimeList manga.
type Manga struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	MainPicture       Picture            `json:"main_picture"`
	AlternativeTitles Titles             `json:"alternative_titles"`
	StartDate         string             `json:"start_date"`
	Synopsis          string             `json:"synopsis"`
	Mean              float64            `json:"mean"`
	Rank              int                `json:"rank"`
	Popularity        int                `json:"popularity"`
	NumListUsers      int                `json:"num_list_users"`
	NumScoringUsers   int                `json:"num_scoring_users"`
	Nsfw              string             `json:"nsfw"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	MediaType         string             `json:"media_type"`
	Status            string             `json:"status"`
	Genres            []Genre            `json:"genres"`
	MyListStatus      MangaListStatus    `json:"my_list_status"`
	NumVolumes        int                `json:"num_volumes"`
	NumChapters       int                `json:"num_chapters"`
	Authors           []Author           `json:"authors"`
	Pictures          []Picture          `json:"pictures"`
	Background        string             `json:"background"`
	RelatedAnime      []RelatedAnime     `json:"related_anime"`
	RelatedManga      []RelatedManga     `json:"related_manga"`
	Recommendations   []RecommendedManga `json:"recommendations"`
	Serialization     []Serialization    `json:"serialization"`
}

// Person is usually the creator of a manga.
type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Author of the manga.
type Author struct {
	Person Person `json:"node"`
	Role   string `json:"role"`
}

// RelatedManga shows manga related with the returned entry.
type RelatedManga struct {
	Node                  Manga  `json:"node,omitempty"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// RecommendedManga is a manga recommended to the user.
type RecommendedManga struct {
	Node               Manga `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// Magazine in which the manga was serialized.
type Magazine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Serialization is a serialized manga series.
type Serialization struct {
	Node Magazine `json:"node"`
	Role string   `json:"role"`
}

// Details returns details about a manga.
func (s *MangaService) Details(ctx context.Context, mangaID int64, options ...DetailsOption) (*Manga, *Response, error) {
	m := new(Manga)
	resp, err := s.client.details(ctx, fmt.Sprintf("manga/%d", mangaID), m, options...)
	if err != nil {
		return nil, resp, err
	}
	return m, resp, nil
}
