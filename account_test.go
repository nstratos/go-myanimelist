package mal

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountService_Verify(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/account/verify_credentials.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		fmt.Fprint(w, `<user><id>1</id><username>TestUser</username></user>`)
	})

	user, _, err := client.Account.Verify()
	if err != nil {
		t.Errorf("Account.Verify returned error: %v", err)
	}

	want := &User{XMLName: xml.Name{Local: "user"}, ID: 1, Username: "TestUser"}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Account.Verify returned %+v, want %+v", user, want)
	}
}

func TestAccountService_Verify_no_content(t *testing.T) {
	setup()
	defer teardown()

	client.SetCredentials("TestUser", "TestPass")
	client.SetUserAgent("TestAgent")

	mux.HandleFunc("/api/account/verify_credentials.xml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBasicAuth(t, r, true, "TestUser", "TestPass")
		testUserAgent(t, r, "TestAgent")
		http.Error(w, "no content", http.StatusNoContent)
	})

	_, _, err := client.Account.Verify()

	if err != NoContentErr {
		t.Errorf("Account.Verify for non existent user expected to return err %v", NoContentErr)
	}
}
