package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AccountService struct {
	client *Client
}

type User struct {
	ID       int    `xml:"id"`
	Username string `xml:"username"`
}

func (s *AccountService) Verify() (*User, *http.Response, error) {
	verifyURL, _ := url.Parse("api/account/verify_credentials.xml")
	u := s.client.BaseURL.ResolveReference(verifyURL)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("User-Agent", s.client.UserAgent)
	req.SetBasicAuth(s.client.Username, s.client.Password)

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	user := User{}
	err = xml.Unmarshal(body, &user)
	if err != nil {
		return nil, resp, fmt.Errorf("cannot unmarshal '%s' (%s)\n", string(body), err)
	}

	return &user, resp, nil
}

// Verify the user's credentials.
//
// Response
//
// Success: 200 status code, XML data for user.
// Failure: 204 status code (no content).
//
// Example Response
//
//  <?xml version="1.0" encoding="utf-8"?>
//  <user>
//    <id>1</id>
//    <username>Xinil</username>
//  </user>
func Verify() (User, error) {
	const verifyURL = "http://myanimelist.net/api/account/verify_credentials.xml"
	req, err := http.NewRequest("GET", verifyURL, nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Add("User-Agent", userAgent)
	req.SetBasicAuth(username, password)

	resp, err := defaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = xml.Unmarshal(body, &user)
	if err != nil {
		return User{}, fmt.Errorf("cannot unmarshal '%s' (%s)\n", string(body), err)
	}

	return user, nil
}
