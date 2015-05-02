package mal

import "encoding/xml"

type AccountService struct {
	client *Client
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

	const verifyURL = "api/account/verify_credentials.xml"

	req, err := s.client.NewRequest("GET", verifyURL, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}
