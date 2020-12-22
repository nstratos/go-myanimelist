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
// JSON decoded and stored in the value pointed to by v. If v implements the
// io.Writer interface, the raw response body will be written to v, without
// attempting to first decode it.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
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
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
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
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	// Re-populate error response body in case JSON unmarshal fails.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	return errorResponse
}

// post sends a POST API request used by Add and Update.
func (c *Client) post(endpoint string, id int, entry interface{}) (*Response, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%d.xml", endpoint, id), entry)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.Background()
	return c.Do(ctx, req, nil)
}

// delete sends a DELETE API request used by Delete.
func (c *Client) delete(endpoint string, id int) (*Response, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", endpoint, id), nil)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return c.Do(ctx, req, nil)
}

// get sends a GET API request used by List and Search.
func (c *Client) get(url string, result interface{}) (*Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return c.Do(ctx, req, result)
}
