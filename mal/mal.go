package mal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Status specifies a status for anime and manga entries.
//type Status int

// Anime and manga entries have a status such as completed, on hold and
// dropped.
//
// Current is for entries marked as currently watching or reading.
//
// Planned is for entries marked as plan to watch or read.
// const (
// 	Current   Status = 1
// 	Completed        = 2
// 	OnHold           = 3
// 	Dropped          = 4
// 	Planned          = 6
// )

const (
	defaultBaseURL             = "https://api.myanimelist.net/v2/"
	defaultListEndpoint        = "malappinfo.php"
	defaultAccountEndpoint     = "api/account/verify_credentials.xml"
	defaultAnimeAddEndpoint    = "anime/"
	defaultAnimeUpdateEndpoint = "api/animelist/update/"
	defaultAnimeDeleteEndpoint = "api/animelist/delete/"
	defaultAnimeSearchEndpoint = "api/anime/search.xml"
	defaultMangaAddEndpoint    = "api/mangalist/add/"
	defaultMangaUpdateEndpoint = "api/mangalist/update/"
	defaultMangaDeleteEndpoint = "api/mangalist/delete/"
	defaultMangaSearchEndpoint = "api/manga/search.xml"
)

// Client manages communication with the MyAnimeList API.
type Client struct {
	client *http.Client

	username string
	password string

	// Base URL for MyAnimeList API requests.
	BaseURL *url.URL

	Account *AccountService
	Anime   *AnimeService
	Manga   *MangaService
}

// Auth is an option that can be passed to NewClient. It allows to specify the
// username and password to be used for authenticating with the MyAnimeList
// API. When this option is used, the client will use basic authentication on
// the requests than need them.
//
// Most API methods require authentication so it is typical to pass this option
// when creating a new client.
func Auth(username, password string) func(*Client) {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// HTTPClient is an option that can be passed to NewClient. It allows to
// specify the HTTP client that will be used to create the requests. If this
// option is not set, a default HTTP client (http.DefaultClient) will be used
// which is usually sufficient.
//
// This option can be set for less trivial cases, when more control over the
// created HTTP requests is required. One such example is providing a timeout
// to cancel requests that exceed it.
func HTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) {
		c.client = httpClient
	}
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
	listEndpoint, _ := url.Parse(defaultListEndpoint)
	accountEndpoint, _ := url.Parse(defaultAccountEndpoint)
	animeAddEndpoint, _ := url.Parse(defaultAnimeAddEndpoint)
	animeUpdateEndpoint, _ := url.Parse(defaultAnimeUpdateEndpoint)
	animeDeleteEndpoint, _ := url.Parse(defaultAnimeDeleteEndpoint)
	animeSearchEndpoint, _ := url.Parse(defaultAnimeSearchEndpoint)
	mangaAddEndpoint, _ := url.Parse(defaultMangaAddEndpoint)
	mangaUpdateEndpoint, _ := url.Parse(defaultMangaUpdateEndpoint)
	mangaDeleteEndpoint, _ := url.Parse(defaultMangaDeleteEndpoint)
	mangaSearchEndpoint, _ := url.Parse(defaultMangaSearchEndpoint)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	c.Account = &AccountService{
		client:   c,
		Endpoint: accountEndpoint,
	}

	c.Anime = &AnimeService{
		client:         c,
		ListEndpoint:   listEndpoint,
		AddEndpoint:    animeAddEndpoint,
		UpdateEndpoint: animeUpdateEndpoint,
		DeleteEndpoint: animeDeleteEndpoint,
		SearchEndpoint: animeSearchEndpoint,
	}

	c.Manga = &MangaService{
		client:         c,
		ListEndpoint:   listEndpoint,
		AddEndpoint:    mangaAddEndpoint,
		UpdateEndpoint: mangaUpdateEndpoint,
		DeleteEndpoint: mangaDeleteEndpoint,
		SearchEndpoint: mangaSearchEndpoint,
	}

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
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// XML decoded and stored in the value pointed to by v. If XML was unable to get
// decoded, it will be returned in Response.Body along with the error so that
// the caller can further inspect it if needed.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req = req.WithContext(ctx)

	dumpRequest(req)
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled, the context's
		// error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()
	dumpResponse(resp)

	response := &Response{Response: resp}
	if err := checkResponse(resp); err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			// if decErr == io.EOF {
			// 	decErr = nil // ignore EOF errors caused by empty response body
			// }
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// ErrNoContent is returned when a MyAnimeList API method returns error 204.
var ErrNoContent = errors.New("no content")

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	//resp := &Response{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot read error response body: %v", err)
	}

	return fmt.Errorf("error response: %q", data)
}

// post sends a POST API request used by Add and Update.
func (c *Client) post(endpoint string, id int, entry interface{}, useAuth bool) (*Response, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%d.xml", endpoint, id), entry)
	if err != nil {
		return nil, err
	}
	if useAuth {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.Background()
	return c.Do(ctx, req, nil)
}

// delete sends a DELETE API request used by Delete.
func (c *Client) delete(endpoint string, id int, useAuth bool) (*Response, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", endpoint, id), nil)
	if err != nil {
		return nil, err
	}
	if useAuth {
		req.SetBasicAuth(c.username, c.password)
	}
	ctx := context.Background()
	return c.Do(ctx, req, nil)
}

// get sends a GET API request used by List and Search.
func (c *Client) get(url string, result interface{}, useAuth bool) (*Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if useAuth {
		req.SetBasicAuth(c.username, c.password)
	}
	ctx := context.Background()
	return c.Do(ctx, req, result)
}
