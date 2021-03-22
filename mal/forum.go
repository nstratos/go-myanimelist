package mal

import (
	"context"
	"fmt"
	"net/url"
	"time"
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

type topicDetail struct {
	Data   TopicDetails `json:"data"`
	Paging Paging       `json:"paging"`
}

func (t topicDetail) pagination() Paging { return t.Paging }

// TopicDetails contain the posts of a forum topic and an optional poll.
type TopicDetails struct {
	Title string `json:"title"`
	Posts []Post `json:"posts"`
	Poll  *Poll  `json:"poll"`
}

// Post is a forum post.
type Post struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy CreatedBy `json:"created_by"`
	Body      string    `json:"body"`
	Signature string    `json:"signature"`
}

// CreatedBy shows the name of the user that created the post or topic.
type CreatedBy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ForumAvator string `json:"forum_avator"`
}

// Poll is an optional poll in a forum post.
type Poll struct {
	ID       int          `json:"id"`
	Question string       `json:"question"`
	Closed   bool         `json:"closed"`
	Options  []PollOption `json:"options"`
}

// PollOption is one of the choices of a poll.
type PollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// A PagingOption includes the Limit and Offset options which are used for
// controlling pagination in results.
type PagingOption interface {
	pagingApply(v *url.Values)
}

// TopicDetails returns details about the forum topic specified by topicID.
func (s *ForumService) TopicDetails(ctx context.Context, topicID int, options ...PagingOption) (TopicDetails, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromPagingOption(options[i])
	}
	d := new(topicDetail)
	resp, err := s.client.list(ctx, fmt.Sprintf("forum/topic/%d", topicID), d, oo...)
	if err != nil {
		return TopicDetails{}, resp, err
	}
	return d.Data, resp, nil
}

func optionFromPagingOption(o PagingOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.pagingApply(v)
	})
}

type topics struct {
	Data   []Topic `json:"data"`
	Paging Paging  `json:"paging"`
}

func (t topics) pagination() Paging { return t.Paging }

// A Topic of the forum.
type Topic struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         CreatedBy `json:"created_by"`
	NumberOfPosts     int       `json:"number_of_posts"`
	LastPostCreatedAt time.Time `json:"last_post_created_at"`
	LastPostCreatedBy CreatedBy `json:"last_post_created_by"`
	IsLocked          bool      `json:"is_locked"`
}

// TopicsOption are options specific to the ForumService.Topics method.
type TopicsOption interface {
	topicsApply(v *url.Values)
}

// BoardID is an option that filters topics based on the board ID.
type BoardID int

func (id BoardID) topicsApply(v *url.Values) { v.Set("board_id", itoa(int(id))) }

// SubboardID is an option that filters topics based on the subboard ID.
type SubboardID int

func (id SubboardID) topicsApply(v *url.Values) { v.Set("subboard_id", itoa(int(id))) }

// sortTopics is an option that sorts the returned topics.
type sortTopics string

// SortTopicsRecent is the default and only sorting value for topics.
const SortTopicsRecent sortTopics = "recent"

func (s sortTopics) topicsApply(v *url.Values) { v.Set("sort", string(s)) }

// Query is an option that allows to search for a term.
type Query string

func (q Query) topicsApply(v *url.Values) { v.Set("q", string(q)) }

// TopicUserName is an option that filters topics based on the topic username.
type TopicUserName string

func (n TopicUserName) topicsApply(v *url.Values) { v.Set("topic_user_name", string(n)) }

// UserName is an option that filters topics based on a username.
type UserName string

func (n UserName) topicsApply(v *url.Values) { v.Set("user_name", string(n)) }

// Topics returns the forum's topics. Make sure to pass at least the Query
// option or you will get an API error.
func (s *ForumService) Topics(ctx context.Context, options ...TopicsOption) ([]Topic, *Response, error) {
	oo := make([]Option, len(options))
	for i := range options {
		oo[i] = optionFromTopicsOption(options[i])
	}
	t := new(topics)
	resp, err := s.client.list(ctx, "forum/topics", t, oo...)
	if err != nil {
		return nil, resp, err
	}
	return t.Data, resp, nil
}

func optionFromTopicsOption(o TopicsOption) optionFunc {
	return optionFunc(func(v *url.Values) {
		o.topicsApply(v)
	})
}
