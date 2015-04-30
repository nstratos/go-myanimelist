package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// MangaEntry holds values such as score, chapter and status that we want our
// manga entry to have when we add or update it on our list.
//
// Status is required and can be:
// 1/reading, 2/completed, 3/onhold, 4/dropped, 6/plantoread
//
// DateStart and DateFinish require 'mmddyyyy' format
//
// EnableDiscussion and EnableRereading can be: 1=enable, 0=disable
//
// Tags are comma separated: test tag, 2nd tag
type MangaEntry struct {
	XMLName            xml.Name `xml:"entry"`
	Volume             int      `xml:"volume,omitempty"`
	Chapter            int      `xml:"chapter,omitempty"`
	Status             string   `xml:"status,omitempty"`
	Score              int      `xml:"score,omitempty"`
	DownloadedChapters int      `xml:"downloaded_chapters,omitempty"`
	TimesReread        int      `xml:"times_reread,omitempty"`
	RereadValue        int      `xml:"reread_value,omitempty"`
	DateStart          string   `xml:"date_start,omitempty"`
	DateFinish         string   `xml:"date_finish,omitempty"`
	Priority           int      `xml:"priority,omitempty"`
	EnableDiscussion   int      `xml:"enable_discussion,omitempty"`
	EnableRereading    int      `xml:"enable_rereading,omitempty"`
	Comments           string   `xml:"comments,omitempty"`
	ScanGroup          string   `xml:"scan_group,omitempty"`
	Tags               string   `xml:"tags,omitempty"`
	RetailVolumes      int      `xml:"retail_volumes,omitempty"`
}

type MangaService struct {
	client *Client
}

func (s *MangaService) Add(mangaID int, entry MangaEntry) (*http.Response, error) {

	const endpoint = "api/mangalist/add/"

	return s.client.post(endpoint, mangaID, entry)
}

func (s *MangaService) Update(mangaID int, entry MangaEntry) (*http.Response, error) {

	const endpoint = "api/mangalist/update/"

	return s.client.post(endpoint, mangaID, entry)
}

func UpdateManga(mangaID int, data MangaEntry) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return nil, fmt.Errorf("cannot marshal: %s", err)
	}

	const updateMangaURL = "http://myanimelist.net/api/mangalist/update/"
	resp, err := postManga(updateMangaURL, mangaID, xmlData)
	if err != nil {
		return nil, fmt.Errorf("update manga failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func AddManga(mangaID int, data MangaEntry) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return nil, fmt.Errorf("cannot marshal: %s", err)
	}

	const addMangaURL = "http://myanimelist.net/api/mangalist/add/"
	resp, err := postManga(addMangaURL, mangaID, xmlData)
	if err != nil {
		return nil, fmt.Errorf("add manga failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func DeleteManga(mangaID int) ([]byte, error) {
	resp, err := deleteManga(mangaID)
	if err != nil {
		return nil, fmt.Errorf("delete manga failed: '%s', %s", string(resp), err)
	}
	fmt.Printf("REPONSE: %v\n", string(resp))

	return resp, nil
}

func postManga(postURL string, mangaID int, data []byte) ([]byte, error) {
	v := url.Values{}
	v.Set("data", string(data))
	fmt.Printf("POST URL: %v\n", fmt.Sprintf("%s%d.xml", postURL, mangaID))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%d.xml", postURL, mangaID), strings.NewReader(v.Encode()))
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

func deleteManga(mangaID int) ([]byte, error) {
	const deleteMangaURL = "http://myanimelist.net/api/mangalist/delete/"
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%d.xml", deleteMangaURL, mangaID), nil)
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
