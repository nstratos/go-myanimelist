package mal

import (
	"net/url"
)

// AccountService handles communication with the account methods of the
// MyAnimeList API.
//
// MyAnimeList API docs: http://myanimelist.net/modules.php?go=api
type AccountService struct {
	client   *Client
	Endpoint *url.URL
}

// User represents a MyAnimeList user. It is returned as a success to Verify.
// type User struct {
// 	XMLName  xml.Name `xml:"user"`
// 	ID       int      `xml:"id"`
// 	Username string   `xml:"username"`
// }

// Verify the user's credentials that the client is using. If verification is
// successful it will return a User with his ID and username. If the verification
// fails it will return an ErrNoContent.
// func (s *AccountService) Verify() (*User, *Response, error) {

// 	user := new(User)
// 	resp, err := s.client.get(s.Endpoint.String(), user, true)
// 	if err != nil {
// 		return nil, resp, err
// 	}
// 	return user, resp, nil
// }
