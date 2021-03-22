package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/nstratos/go-myanimelist/mal"
	"golang.org/x/oauth2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// In order to create a client ID and secret for your application:
//
//  1. Navigate to https://myanimelist.net/apiconfig or go to your MyAnimeList
//     profile, click Edit Profile and select the API tab on the far right.
//  2. Click Create ID and submit the form with your application details.
const (
	defaultClientID     = ""
	defaultClientSecret = ""
)

// Authorization Documentation:
//
// https://myanimelist.net/apiconfig/references/authorization

func run() error {
	var (
		clientID     = flag.String("client-id", defaultClientID, "your application client ID")
		clientSecret = flag.String("client-secret", defaultClientSecret, "your application client secret")
		// state is a token to protect the user from CSRF attacks. In a web
		// application, you should provide a non-empty string and validate that
		// it matches the state query parameter on the redirect URL callback
		// after the MyAnimeList authentication. It can stay empty here.
		state = flag.String("state", "", "token to protect against CSRF attacks")
	)
	flag.Parse()

	ctx := context.Background()

	tokenClient, err := authenticate(ctx, *clientID, *clientSecret, *state)
	if err != nil {
		return err
	}

	c := demoClient{
		Client: mal.NewClient(tokenClient),
	}

	return c.showcase(ctx)
}

func authenticate(ctx context.Context, clientID, clientSecret, state string) (*http.Client, error) {
	accessToken := loadCachedToken()
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		return oauth2.NewClient(ctx, ts), nil
	}

	// Prepare the oauth2 configuration with your application ID, secret, the
	// MyAnimeList authentication and token URLs as specified in:
	//
	// https://myanimelist.net/apiconfig/references/authorization
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	// Generate a code verifier, a high-entropy cryptographic random string. It
	// will be set as the code_challenge in the authentication URL. It should
	// have a minimum length of 43 characters and a maximum length of 128
	// characters.
	const codeVerifierLength = 128
	codeVerifier, err := generateCodeVerifier(codeVerifierLength)
	if err != nil {
		return nil, fmt.Errorf("generating code verifier: %v", err)
	}

	// Produce the authentication URL where the user needs to be redirected and
	// allow your application to access their MyAnimeList data.
	authURL := conf.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeVerifier),
	)
	err = openBrowser(authURL)
	if err != nil {
		fmt.Println("Could not open browser.")
	}

	fmt.Printf("Your browser should open: %v\n", authURL)
	fmt.Print("After authenticating, copy the code from the browser URL and paste it here: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	code := scanner.Text()
	if scanner.Err() != nil {
		return nil, fmt.Errorf("reading code from terminal: %v", err)
	}

	// Exchange the authentication code for a token. MyAnimeList currently only
	// supports the plain code_challenge_method so to verify the string, just
	// make sure it is the same as the one you entered in the code_challenge.
	token, err := conf.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %v", err)
	}
	fmt.Println("Authentication was successful. Caching access token...")
	cacheToken(token.AccessToken)

	return conf.Client(ctx, token), nil
}

const cacheName = "auth-example-token-cache.txt"

func cacheToken(token string) {
	content := []byte(token)
	err := os.WriteFile(cacheName, content, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "caching access token: %v, token is: %s", err, token)
		return
	}
}

func loadCachedToken() string {
	token, err := os.ReadFile(cacheName)
	if err != nil {
		return ""
	}
	return string(token)
}

func generateCodeVerifier(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstvuwxyz" +
		"0123456789-._~"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("openBrowser: unsupported operating system: %v", runtime.GOOS)
	}
}
