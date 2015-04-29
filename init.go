package mal

import (
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
