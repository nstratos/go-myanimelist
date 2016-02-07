package mal

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "http://myanimelist.net/"

	defaultUserAgent = `
	Mozilla/5.0 (X11; Linux x86_64) 
	AppleWebKit/537.36 (KHTML, like Gecko) 
	Chrome/42.0.2311.90 Safari/537.36`

	defaultListEndpoint        = "malappinfo.php"
	defaultAccountEndpoint     = "api/account/verify_credentials.xml"
	defaultAnimeAddEndpoint    = "api/animelist/add/"
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

	// User agent used when communicating with the MyAnimeList API.
	UserAgent string
	Username  string
	Password  string

	// Base URL for MyAnimeList API requests.
	BaseURL *url.URL

	Account *AccountService
	Anime   *AnimeService
	Manga   *MangaService
}

// NewClient returns a new MyAnimeList API client.
func NewClient() *Client {
	httpClient := http.DefaultClient

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
		client:    httpClient,
		UserAgent: defaultUserAgent,
		BaseURL:   baseURL,
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

// SetCredentials sets the username and password that will be used for basic
// authentication.
func (c *Client) SetCredentials(username, password string) {
	c.Username = username
	c.Password = password
}

// SetUserAgent sets the user agent that will be used to communicate with the
// MyAnimeList API. If no user agent is provided then a default one will be used.
//
// MyAnimeList uses the user agent as a token to identify applications. It is
// important to get your own whitelisted user agent if you are planning to use
// this library in your application. Otherwise your IP might get blocked due to
// excessive requests.
//
// To get your own whitelisted user agent, see:
// http://myanimelist.net/forum/?topicid=692311
func (c *Client) SetUserAgent(userAgent string) {
	c.UserAgent = userAgent
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
func (c *Client) NewRequest(method, urlStr string, data interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	v := url.Values{}
	if data != nil {
		d, err := xml.Marshal(data)
		if err != nil {
			return nil, err
		}
		v.Set("data", string(d))
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return req, nil

}

// Do sends an API request and returns the API response. The API response is
// XML decoded and stored in the value pointed to by v. If XML was unable to get
// decoded, it will be returned in Response.Body along with the error so that
// the caller can further inspect it if needed.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := readResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		b := response.Body
		// enconding/xml cannot handle entity &bull;
		b = bytes.Replace(b, []byte("&bull;"), []byte("<![CDATA[&bull;]]>"), -1)
		err := xml.Unmarshal(b, v)
		if err != nil {
			return response, fmt.Errorf("cannot decode: %v", err)
		}
	}

	return response, nil
}

// ErrNoContent is returned when a MyAnimeList API method returns error 204.
var ErrNoContent = errors.New("no content")

func readResponse(r *http.Response) (*Response, error) {
	resp := &Response{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return resp, fmt.Errorf("cannot read response body: %v", err)
	}
	resp.Body = data

	if r.StatusCode == http.StatusNoContent {
		return resp, ErrNoContent
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return resp, fmt.Errorf("%v %v: %d %s",
			r.Request.Method, r.Request.URL,
			r.StatusCode, string(data))
	}

	return resp, nil
}

// post sends a POST API request used by Add and Update.
func (c *Client) post(endpoint string, id int, entry interface{}) (*Response, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%d.xml", endpoint, id), entry)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// post sends a DELETE API request used by Delete.
func (c *Client) delete(endpoint string, id int) (*Response, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", endpoint, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// get sends a GET API request used by List and Search.
func (c *Client) get(url string, result interface{}) (*Response, error) {

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, result)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
