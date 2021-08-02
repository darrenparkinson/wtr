package auth

import (
	"context"
	"crypto/rand"
	_ "embed" // for embedding success and error html
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	//go:embed success.html
	successHTML string
	//go:embed error.html
	errorHTML string
)

const (
	// AuthURL is the URL to Webex Accounts Service's OAuth2 endpoint.
	AuthURL = "https://webexapis.com/v1/authorize"
	// TokenURL is the URL to the Webex Accounts Service's OAuth2
	// token endpoint.
	TokenURL = "https://webexapis.com/v1/access_token"
)

// Webex is a top level struct for our communication with Webex
type Webex struct {
	Client       *http.Client // used for making request to webex
	ClientID     string       // Webex App Integration ClientID
	ClientSecret string       // Webex App Integration Client Secret
	RedirectPort string       // port on localhost to listen for webex response, as configured in the app integration redirectURI as http://localhost:PORT
	Scopes       []string     // scopes configured on Webex App Integration
}

// GetAccessToken retrieves an integration token from Webex.  It starts a server on
// localhost at the port specified, sends the user off for authentication and then
// retrieves an WebexAccessTokenResponse with the returned code from the authorization code
// grant for later use.  Default listener is http://localhost:6855.
func (wbx *Webex) GetAccessToken(timeout int) (*WebexAccessTokenResponse, error) {
	// Check we have required config
	if wbx.ClientID == "" || wbx.ClientSecret == "" {
		return nil, ErrMissingCredentials
	}
	if wbx.RedirectPort == "" {
		wbx.RedirectPort = "6855"
	}
	if wbx.Scopes == nil {
		wbx.Scopes = []string{"spark:kms", "spark:all"}
	}

	state := generateRandomState()

	srv, done, err := startServer(wbx, state)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server started: listening for webex response")

	authURL := wbx.AuthURL(state)
	openBrowser(authURL)

	token, err := waitForResponse(done, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("response received")

	log.Println("stopping server")
	err = stopServer(srv)
	if err != nil {
		log.Println("error stopping server:", err)
	}
	log.Println("server stopped")

	return token, nil
}

// RefreshToken refreshes an existing token given the appropriate details.
func RefreshToken(clientID, secret, refreshToken string) (*WebexAccessTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("refresh_token", refreshToken)
	return makeTokenRequest(TokenURL, form)
}

// AuthURL builds the required URL to send the user given the values in Webex
func (wbx *Webex) AuthURL(state string) string {
	authURL, _ := url.Parse(AuthURL)
	q := authURL.Query()
	q.Set("client_id", wbx.ClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", fmt.Sprintf("http://localhost:%s", wbx.RedirectPort))
	q.Set("scope", strings.Join(wbx.Scopes, " "))
	q.Set("state", state)
	authURL.RawQuery = q.Encode()
	return authURL.String()
}

func startServer(wbx *Webex, state string) (*http.Server, chan *WebexAccessTokenResponse, error) {
	done := make(chan *WebexAccessTokenResponse)
	router := mux.NewRouter()
	router.HandleFunc("/", waitHandler(done, state, wbx)).Methods("GET")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", wbx.RedirectPort),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv, done, nil
}

func waitForResponse(done chan *WebexAccessTokenResponse, timeout time.Duration) (*WebexAccessTokenResponse, error) {
	select {
	case token := <-done:
		return token, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out waiting for authentication")
	}
}

func stopServer(srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("server shutdown failed")
	}
	return nil
}

func waitHandler(done chan *WebexAccessTokenResponse, state string, wbx *Webex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := extractCode(state, r)
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, errorHTML)
			return
		}
		token, err := exchangeCodeForToken(code, wbx)
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, errorHTML)
			return
		}
		fmt.Fprintln(w, successHTML)
		done <- token
	}
}

func extractCode(state string, r *http.Request) (string, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return "", errors.New("webex: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return "", errors.New("webex: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return "", errors.New("webex: redirect state parameter doesn't match")
	}
	return code, nil
}

func exchangeCodeForToken(code string, wbx *Webex) (*WebexAccessTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", wbx.ClientID)
	form.Set("client_secret", wbx.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", fmt.Sprintf("http://localhost:%s", wbx.RedirectPort))
	return makeTokenRequest(TokenURL, form)
}

func makeTokenRequest(tokenURL string, form url.Values) (*WebexAccessTokenResponse, error) {
	client := &http.Client{}
	var tr WebexAccessTokenResponse

	r, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return nil, err
	}
	if tr.Message != "" {
		return nil, errors.New(tr.Message)
	}
	return &tr, nil
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func generateRandomState() string {
	buf := make([]byte, 7)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	state := hex.EncodeToString(buf)
	return state
}
