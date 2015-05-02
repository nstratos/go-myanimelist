package mal

import (
	"bytes"
	"encoding/xml"
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

	defaultListURL        = "malappinfo.php"
	deleteMangaURL        = "api/mangalist/"
	defaultAnimeURL       = "api/animelist/"
	defaultAccountURL     = "api/account/verify_credentials.xml"
	defaultSearchAnimeURL = "api/anime/search.xml"
	defaultSearchMangaURL = "api/manga/search.xml"
)

type Client struct {
	client *http.Client

	// User agent used when communicateing with the myAnimeList API.
	UserAgent string
	Username  string
	Password  string

	// BaseURL for myAnimeList API requests.
	BaseURL *url.URL

	Account *AccountService
	Anime   *AnimeService
	Manga   *MangaService
}

func NewClient() *Client {
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: defaultUserAgent}
	c.Account = &AccountService{client: c}
	c.Anime = &AnimeService{client: c}
	c.Manga = &MangaService{client: c}
	return c
}

func (c *Client) SetCredentials(username, password string) {
	c.Username = username
	c.Password = password
}

func (c *Client) SetUserAgent(userAgent string) {
	c.UserAgent = userAgent
}

type Response struct {
	*http.Response
	Body []byte
}

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

	//if v != nil && len(response.Body) != 0 {
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

func readResponse(r *http.Response) (*Response, error) {
	resp := &Response{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return resp, fmt.Errorf("cannot read response body: %v", err)
	}
	resp.Body = data

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return resp, fmt.Errorf("%v %v: %d %s",
			r.Request.Method, r.Request.URL,
			r.StatusCode, string(data))
	}

	return resp, nil
}

func (c *Client) post(endpoint string, id int, entry interface{}) (*Response, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%d.xml", endpoint, id), entry)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) delete(endpoint string, id int) (*Response, error) {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", endpoint, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) query(url string, result interface{}) (*Response, error) {

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
