package mal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultBaseURL = "https://api.myanimelist.net/v2/"
)

// Client manages communication with the MyAnimeList API.
type Client struct {
	client *http.Client

	// Base URL for MyAnimeList API requests.
	BaseURL *url.URL

	Anime *AnimeService
	Manga *MangaService
	User  *UserService
	Forum *ForumService
}

// NewClient returns a new MyAnimeList API client. The httpClient parameter
// allows to specify the http.client that will be used for all API requests. If
// a nil httpClient is provided, a new http.Client will be used.
//
// In the typical case, you will want to provide an http.Client that will
// perform the authentication for you. Such a client is provided by the
// golang.org/x/oauth2 package. Check out the example directory of the project
// for a full authentication example.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	c.User = &UserService{client: c}
	c.Anime = &AnimeService{client: c}
	c.Manga = &MangaService{client: c}
	c.Forum = &ForumService{client: c}

	return c
}

// Response wraps http.Response and is returned in all the library functions
// that communicate with the MyAnimeList API. Even if an error occurs the
// response will always be returned along with the actual error so that the
// caller can further inspect it if needed. For the same reason it also keeps
// a copy of the http.Response.Body that was read when the response was first
// received.
type Response struct {
	*http.Response
	Body []byte

	NextOffset int
	PrevOffset int
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If data
// is passed as an argument then it will first be encoded in XML and then added
// to the request body as URL encoded value data=<xml>...
// This is how the MyAnimeList requires to receive the data when adding or
// updating entries.
//
// MyAnimeList API docs: http://myanimelist.net/modules.php?go=api
func (c *Client) NewRequest(method, urlStr string, urlOptions ...func(v *url.Values)) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var r io.Reader
	if len(urlOptions) != 0 {
		v := &url.Values{}
		for _, o := range urlOptions {
			o(v)
		}
		r = strings.NewReader(v.Encode())
	}

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	if len(urlOptions) != 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v. If v implements the
// io.Writer interface, the raw response body will be written to v, without
// attempting to first decode it.
//
// If the provided ctx is nil then an error will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	req = req.WithContext(ctx)

	dumpRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dumpResponse(resp)

	response := &Response{Response: resp}
	if err := checkResponse(resp); err != nil {
		return response, err
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}

	return response, err
}

// An ErrorResponse reports an error caused by an API request.
//
// https://myanimelist.net/apiconfig/references/api/v2#section/Common-formats
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"`
	Err      string         `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Err)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		// Ignore unmarshal error for undocumented error formats or HTML.
		_ = json.Unmarshal(data, errorResponse)
	}
	// Re-populate error response body in case JSON unmarshal fails.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	return errorResponse
}

func (c *Client) details(ctx context.Context, path string, v interface{}, options ...DetailsOption) (*Response, error) {
	req, err := c.NewRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for _, o := range options {
		o.detailsApply(&q)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(ctx, req, v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Paging provides access to the next and previous page URLs when there are
// pages of results.
type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type pagination interface {
	pagination() Paging
}

func (c *Client) list(ctx context.Context, path string, p pagination, options ...Option) (*Response, error) {
	req, err := c.NewRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for _, o := range options {
		o.apply(&q)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(ctx, req, p)
	if err != nil {
		return resp, err
	}

	prev, next, err := parsePaging(p.pagination())
	if err != nil {
		return resp, err
	}
	resp.PrevOffset = prev
	resp.NextOffset = next

	return resp, nil
}

func parsePaging(p Paging) (prev, next int, err error) {
	if p.Previous != "" {
		offset, err := parseOffset(p.Previous)
		if err != nil {
			return 0, 0, fmt.Errorf("paging: previous: %s", err)
		}
		prev = offset
	}
	if p.Next != "" {
		offset, err := parseOffset(p.Next)
		if err != nil {
			return 0, 0, fmt.Errorf("paging: next: %s", err)
		}
		next = offset
	}
	return prev, next, nil
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
