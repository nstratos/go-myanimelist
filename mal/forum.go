package mal

import (
	"context"
)

// ForumService handles communication with the forum related methods of the
// MyAnimeList API:
//
// https://myanimelist.net/apiconfig/references/api/v2#tag/forum
type ForumService struct {
	client *Client
}

// The Forum of MyAnimeList.
type Forum struct {
	Categories []ForumCategory `json:"categories"`
}

// ForumCategory is a category of the forum.
type ForumCategory struct {
	Title  string       `json:"title"`
	Boards []ForumBoard `json:"boards"`
}

// ForumBoard is a board of the forum.
type ForumBoard struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Subboards   []ForumSubboard `json:"subboards"`
}

// ForumSubboard is a subboard of the forum.
type ForumSubboard struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Boards returns the forum boards.
func (s *ForumService) Boards(ctx context.Context) (*Forum, *Response, error) {
	f := new(Forum)
	resp, err := s.client.details(ctx, "forum/boards", f)
	if err != nil {
		return nil, resp, err
	}
	return f, resp, nil
}
