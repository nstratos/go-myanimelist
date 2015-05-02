package mal

import (
	"encoding/xml"
	"net/url"
)

type AccountService struct {
	client   *Client
	Endpoint *url.URL
}

type User struct {
	XMLName  xml.Name `xml:"user"`
	ID       int      `xml:"id"`
	Username string   `xml:"username"`
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
func (s *AccountService) Verify() (*User, *Response, error) {

	user := new(User)
	resp, err := s.client.query(s.Endpoint.String(), user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}
