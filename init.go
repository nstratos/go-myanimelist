package mal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "http://myanimelist.net/"

	defaultUserAgent = "api-indiv-2D4068FCF43349DA30D8D4E5667883C2"

	animeListURL = "http://myanimelist.net/malappinfo.php?status=all&type=anime&u="
	mangaListURL = "http://myanimelist.net/malappinfo.php?status=all&type=manga&u="

	updateMangaURL = "http://myanimelist.net/api/mangalist/update/"
	addMangaURL    = "http://myanimelist.net/api/mangalist/add/"
	deleteMangaURL = "http://myanimelist.net/api/animelist/delete/"

	updateAnimeURL = "http://myanimelist.net/api/animelist/update/"
	addAnimeURL    = "http://myanimelist.net/api/animelist/add/"
	deleteAnimeURL = "http://myanimelist.net/api/animelist/delete/"

	searchAnimeURL = "http://myanimelist.net/api/anime/search.xml?q="
	searchMangaURL = "http://myanimelist.net/api/manga/search.xml?q="

	verifyURL = "http://myanimelist.net/api/account/verify_credentials.xml"
)

var defaultClient = &http.Client{}

var username, password, userAgent string

func Init(uname, passwd, agent string) {
	username = uname
	password = passwd
	userAgent = agent
}

type Client struct {
	client *http.Client

	// User agent used when communicateing with the myAnimeList API.
	UserAgent string
	Username  string
	Password  string

	// BaseURL for myAnimeList API requests.
	BaseURL *url.URL

	Account *AccountService
}

func NewClient() *Client {
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: defaultUserAgent}
	c.Account = &AccountService{client: c}
	return c
}

func (c *Client) SetCredentials(username, password, userAgent string) {
	c.Username = username
	c.Password = password
	c.UserAgent = userAgent
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := xml.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
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

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		// If we get a writer, we copy the body to the writer (e.g. so it can
		// write in a created file) otherwise we decode the expected XML.
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = xml.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body %v", err)
	}
	return &ErrorResponse{Response: r, Body: data}
}

type ErrorResponse struct {
	Response *http.Response
	Body     []byte
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d '%s'",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, string(r.Body))
}
