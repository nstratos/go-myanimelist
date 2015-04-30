package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AnimeService struct {
	client *Client
}

func (s *AnimeService) Add(animeID int, data AnimeData) (*http.Response, error) {

	const endpoint = "api/animelist/add/"

	req, err := s.client.NewRequest("POST", fmt.Sprintf("%s%d.xml", endpoint, animeID), data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AnimeService) Delete(animeID int) (*http.Response, error) {

	const endpoint = "api/animelist/delete/"

	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", endpoint, animeID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AnimeData holds values such as score, episode and status that we want our
// anime entry to have when we add or update it on our list.
//
// Status is required and can be:
// 1/watching, 2/completed, 3/onhold, 4/dropped, 6/plantowatch
//
// DateStart and DateFinish require 'mmddyyyy' format
//
// EnableDiscussion and EnableRewatching can be: 1=enable, 0=disable
//
// Tags are comma separated: test tag, 2nd tag
type AnimeData struct {
	XMLName            xml.Name `xml:"entry"`
	Episode            int      `xml:"episode,omitempty"`
	Status             string   `xml:"status,omitempty"`
	Score              int      `xml:"score,omitempty"`
	DownloadedEpisodes int      `xml:"downloaded_episodes,omitempty"`
	StorageType        int      `xml:"storage_type,omitempty"`
	StorageValue       float64  `xml:"storage_value,omitempty"`
	TimesRewatched     int      `xml:"times_rewatched,omitempty"`
	RewatchValue       int      `xml:"rewatch_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`
	DateFinish         string   `xml:"date_finish,omitempty"`
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"`
	EnableRewatching   int      `xml:"enable_rewatching,omitempty"`
	Comments           string   `xml:"comments,omitempty"`
	FansubGroup        string   `xml:"fansub_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"`
}

func UpdateAnime(animeID int, data AnimeData) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return nil, fmt.Errorf("cannot marshal: %s", err)
	}

	const updateAnimeURL = "http://myanimelist.net/api/animelist/update/"
	resp, err := postAnime(updateAnimeURL, animeID, xmlData)
	if err != nil {
		return nil, fmt.Errorf("update anime failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func AddAnime(animeID int, data AnimeData) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return nil, fmt.Errorf("cannot marshal: %s", err)
	}

	const addAnimeURL = "http://myanimelist.net/api/animelist/add/"
	resp, err := postAnime(addAnimeURL, animeID, xmlData)
	if err != nil {
		return nil, fmt.Errorf("add anime failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func DeleteAnime(animeID int) ([]byte, error) {
	resp, err := deleteAnime(animeID)
	if err != nil {
		return nil, fmt.Errorf("delete anime failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func postAnime(postURL string, animeID int, data []byte) ([]byte, error) {
	v := url.Values{}
	v.Set("data", string(data))
	fmt.Printf("POST URL: %v\n", fmt.Sprintf("%s%d.xml", postURL, animeID))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%d.xml", postURL, animeID), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func deleteAnime(animeID int) ([]byte, error) {
	const deleteAnimeURL = "http://myanimelist.net/api/animelist/delete/"
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", deleteAnimeURL, animeID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.SetBasicAuth(username, password)

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
