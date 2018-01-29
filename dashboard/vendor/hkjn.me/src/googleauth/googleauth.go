// Package googleauth provides OAuth sign-in using Google+.
//
// Minimum setup needed:
//   1. Specify client G+ credentials with SetCredentials()
//   2. Specify a gating function with SetGatingFunc()
//   3. Register a HTTP GET route at /connect for ConnectHandler
//   4. Wrap any HTTP routes that should be authenticated with RequireLogin()
//
// This will send users accessing the resources under authentication
// to a simple page with a G+ button, and if SetGatingFunc accepts
// that G+ id, the user is redirected to the original URL.
//
// The login page can be changed by setting a different LoginTmpl.
//
// If more control is desired, IsLoggedIn, LogIn and Connect can be
// used directly, but with the steps above it's not necessary.
package googleauth

// TODO: abstract away use of glog in favor of a generic logger; it
// would break on app engine since it writes to local disk.

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	// Default policy is to deny all users. Use SetGatingFunc.
	isAllowed   = func(string) bool { return false }
	sessionName = "gplusclient"                                     // name of Gorilla session cookie
	store       = sessions.NewCookieStore([]byte(randomString(32))) // Gorilla session store
	oauthConfig = &oauth2.Config{                                   // config supplied to the OAuth package
		// Scope determines which API calls we are authorized to make. We
		// only want the basics.
		Scopes: []string{"profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		// Setting "postmessage" means that we're handling the access
		// token exchange server-side.
		RedirectURL: "postmessage",
	}
	accessDeniedErr = errors.New("user is not allowed access")
)

// isAccessDenied returns whether the error is access denied.
func isAccessDenied(err error) bool {
	return err.Error() == accessDeniedErr.Error()
}

// SetCredential specifies the client G+ credentials.
func SetCredentials(clientId, clientSecret string) {
	oauthConfig.ClientID = clientId
	oauthConfig.ClientSecret = clientSecret
}

// SetGatingFunc sets a function to check if user with given G+ id is
// allowed access.
func SetGatingFunc(fn func(gplusId string) bool) {
	isAllowed = fn
}

// IsLoggedIn returns true if the user is signed in.
func IsLoggedIn(r *http.Request) (bool, error) {
	session, err := getSession(r)
	if err != nil {
		return false, err
	}
	t := session.Values["accessToken"]
	if t == nil {
		return false, nil
	}
	storedToken, ok := t.(string)
	if !ok {
		return false, fmt.Errorf("bad type of %q value in session: %v", "accessToken", err)
	}
	gp := session.Values["gplusID"]
	if t == nil {
		return false, nil
	}
	gplusId, ok := gp.(string)
	if !ok {
		return false, fmt.Errorf("bad type of %q value in session: %v", "gplusID", err)
	}
	return storedToken != "" && isAllowed(gplusId), nil
}

// LoginInfo represents the user's login info.
type LoginInfo struct {
	ClientId   string // id of client
	StateToken string // state token
}

// LogIn returns the user's login info, starting the auth process.
//
// LogIn generates a state token, which along with the client id
// should be returned to the user where the front-end library can
// exchange them for a one-time authorization code. That one-time
// authorization code is then passed in to /connect, which finishes
// the auth process.
func LogIn(w http.ResponseWriter, r *http.Request) (*LoginInfo, error) {
	// Create a state token to prevent request forgery and store it in the session
	// for later validation.
	session, err := getSession(r)
	if err != nil {
		return nil, err
	}

	state := randomString(64)
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		return nil, fmt.Errorf("failed to save state in session: %v", err)
	}
	glog.V(1).Infof("CheckLogin set state=%v in user's session\n", state)
	return &LoginInfo{oauthConfig.ClientID, url.QueryEscape(state)}, nil
}

// Connect finishes the connection process, exchanging the one-time
// authorization code for an access token and storing it in the
// session.
func Connect(w http.ResponseWriter, r *http.Request) error {
	loggedIn, err := IsLoggedIn(r)
	if err != nil {
		return err
	}
	if loggedIn {
		glog.V(1).Infof("already logged in\n")
		return nil
	}

	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user.
	session, err := getSession(r)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	if session.Values["state"] == nil {
		return fmt.Errorf("missing %q variable in session for user trying to log in? bug, or user is trying to spoof log in", "state")
	}
	sessionState := session.Values["state"].(string)
	if q.Get("state") != sessionState {
		// Note: This can happen if CheckLogIn is called multiple times
		// for the same session, e.g. when several tabs are loading
		// protected resources.
		return fmt.Errorf("state mismatch, got %q from form, but had %q in session\n", r.FormValue("state"), sessionState)
	}
	session.Values["state"] = nil

	code := q.Get("code")
	if code == "" {
		return fmt.Errorf("missing %q value in request body", "code")
	}
	glog.V(1).Infof("code=%v\n", code)
	// We got back matching state from user as well as auth code from
	// login button, exchange the one-time auth code for access token +
	// user id.
	accessToken, idToken, err := exchange(code)
	if err != nil {
		return fmt.Errorf("couldn't exchange code for access token: %v", err)
	}
	glog.V(1).Infof("id token: %v\n", idToken)
	gplusId, err := decodeIdToken(idToken)
	glog.V(1).Infof("decoded G+ token: %v\n", gplusId)
	if err != nil {
		return fmt.Errorf("couldn't decode ID token: %v", err)
	}

	if !isAllowed(gplusId) {
		glog.Infof("user with G+ %v is not allowed access\n", gplusId)
		return accessDeniedErr
	}
	glog.V(1).Infof("User %v is allowed to log in\n", gplusId)

	// Store the access token in the session for later use.
	session.Values["accessToken"] = accessToken
	session.Values["gplusID"] = gplusId
	err = session.Save(r, w)
	if err != nil {
		return fmt.Errorf("failed to save state in session: %v", err)
	}
	return nil
}

// exchange takes an authentication code and exchanges it with the OAuth
// endpoint for a Google API bearer token and a Google+ ID.
func exchange(code string) (accessToken string, idToken string, err error) {
	if oauthConfig.ClientID == "" {
		return "", "", fmt.Errorf("missing client id")
	}
	if oauthConfig.ClientSecret == "" {
		return "", "", fmt.Errorf("missing client secret")
	}
	// Exchange the authorization code for a credentials object via a
	// POST request.
	addr := "https://accounts.google.com/o/oauth2/token"
	values := url.Values{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"code":          {code},
		"client_id":     {oauthConfig.ClientID},
		"client_secret": {oauthConfig.ClientSecret},
		"redirect_uri":  {oauthConfig.RedirectURL},
		"grant_type":    {"authorization_code"},
	}
	resp, err := http.PostForm(addr, values)
	if err != nil {
		return "", "", fmt.Errorf("error exchanging code: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response body into a token object.
	token := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		IdToken     string `json:"id_token"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", "", fmt.Errorf("error decoding access token: %v", err)
	}

	return token.AccessToken, token.IdToken, nil
}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within.
func decodeIdToken(idToken string) (gplusID string, err error) {
	// An ID token is a cryptographically-signed JSON object encoded in
	// base 64.  Normally, it is critical to validate an ID token before
	// using it, but since we are communicating directly with Google
	// over an intermediary-free HTTPS channel and using the Client
	// Secret to authenticate ourselves, we can be confident that the
	// token we receive really comes from Google and is valid. If this
	// is ever passed outside the googleauth package, it is extremely
	// important to validate the token before using it.
	set := struct{ Sub string }{}
	if idToken != "" {
		// Check that the padding is correct for a base64decode.
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("bad ID token")
		}
		// Decode the ID token.
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("bad ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("bad ID token: %v", err)
		}
	}
	return set.Sub, nil
}

// getSession returns the user's Gorilla session.
func getSession(r *http.Request) (*sessions.Session, error) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		if session.IsNew {
			glog.V(1).Infof("ignoring initial session fetch error since session IsNew: %v\n", err)
		} else {
			return nil, fmt.Errorf("error fetching session: %v", err)
		}
	}
	return session, nil
}

// randomString returns a random string with the specified length
func randomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// base64Decode decodes specified base64 string.
func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
