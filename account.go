package mal

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	ID       int    `xml:"id"`
	Username string `xml:"username"`
}

// Response:
// Success: 200 status code, XML data for user.
// Failure: 204 status code (no content).

// Example Response:
// <?xml version="1.0" encoding="utf-8"?>
// <user>
//     <id>1</id>
//     <username>Xinil</username>
// </user>
func Verify() (User, error) {
	const verifyUrl = "http://myanimelist.net/api/account/verify_credentials.xml"
	req, err := http.NewRequest("GET", verifyUrl, nil)
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
