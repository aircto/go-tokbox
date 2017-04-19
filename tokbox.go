package tokbox

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/twinj/uuid"
)

const (
	apiURL = "https://api.opentok.com"
)

// Tokbox represents a main tokbox type that wraps the REST API.
type Tokbox struct {
	key    string
	secret string
}

// New is a constructor function that takes Key and Secret and returns
// a tokbox instance.
func New(key, secret string) *Tokbox {
	return &Tokbox{
		key:    key,
		secret: secret,
	}
}

// CreateSession creates a unique session, to which clients can connect to.
// It takes nothing and return a session instance.
func (t *Tokbox) CreateSession() (*Session, error) {
	url := apiURL + "/session/create"
	token, err := t.Token()
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"X-OPENTOK-AUTH": token,
	}
	res, err := t.MakeRequest("POST", url, headers, map[string]string{
		"archiveMode": "always",
	})
	if err != nil {
		return nil, errors.Wrap(err, "newrequest")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("create session failed " + res.Status)
	}

	var s []Session // must be a list according to docs. https://tokbox.com/developer/rest/#session_id_production
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, errors.Wrap(err, "decoding")
	}
	if len(s) != 1 {
		return nil, errors.New("api returned more than 1 responses.")
	}
	return &s[0], nil
}

// jwtToken generates unique jwt token every time its called.
// Its used make any api request to tokbox REST API.
func (t *Tokbox) Token() (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    t.key,
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Unix() + (2 * 24 * 60 * 60), // 2 hours
		Id:        string(uuid.NewV4()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secret))
}

// Archives returns a list of archived media for a given sessionID.
func (t *Tokbox) Archives(sessionID string) ([]Archive, error) {
	url := apiURL + "/v2/project/" + t.key + "/archive?sessionId=" + sessionID
	token, err := t.Token()
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"X-OPENTOK-AUTH": token,
	}
	res, err := t.MakeRequest("GET", url, headers, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("retrieve archive list failed. " + res.Status)
	}
	jResp := struct {
		Count int       `json:"count"`
		Items []Archive `json:"items"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&jResp); err != nil {
		return nil, errors.Wrap(err, "archive decode")
	}
	return jResp.Items, nil
}

// NewRequest create a single http.Request based on url, headers and body that needs
// to be encoded. Returns the http.Response. any non-nil error means request is
// unsuccessfull.
func (t *Tokbox) MakeRequest(
	method, urlStr string,
	headers map[string]string,
	body interface{},
) (*http.Response, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, rel.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return http.DefaultClient.Do(req)
}
