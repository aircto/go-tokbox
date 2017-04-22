package tokbox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/twinj/uuid"
)

const (
	baseURL       = "https://api.opentok.com"
	tokenSentinel = "T1=="
	apiVersion    = "v2"
)

// ArchiveMode denotes different modes for an archive as per doc.
// https://tokbox.com/developer/rest/#start_archive
type ArchiveMode string

const (
	ArchiveModeManual ArchiveMode = "manual"
	ArchiveModeAlways ArchiveMode = "always"
)

// ArchiveOutputMode lists a valid output mode for an archive
type ArchiveOutputMode string

const (
	ArchiveOutputModeComposed   ArchiveOutputMode = "composed"
	ArchiveOutputModeIndividual ArchiveOutputMode = "individual"
)

// Role is a list of valid token for an token
type Role string

const (
	RoleModerator  string = "moderator"
	RolePublisher  string = "publisher"
	RoleSubscriber string = "subscriber"
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
func (t *Tokbox) CreateSession() (Session, error) {
	url := baseURL + "/session/create"

	// NOTE(kaviraj): According to create session doc. Tokbox returns list of sessions
	// even it creates just one session
	var sessions []Session
	if err := t.MakeRequest("POST", url, map[string]string{
		"archiveMode": "manual",
	}, &sessions); err != nil {
		return Session{}, err
	}

	if len(sessions) != 1 {
		return Session{}, errors.New("api returned more than 1 responses.")
	}
	return sessions[0], nil
}

// jwtToken generates unique jwt token every time its called.
// Its used make any api request to tokbox REST API.
func (t *Tokbox) jwtToken() (string, error) {
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
	url := baseURL + "/" + apiVersion + "/project/" + t.key + "/archive?sessionId=" + sessionID

	var archives ArchiveList
	if err := t.MakeRequest("GET", url, nil, &archives); err != nil {
		return nil, err
	}

	return archives.Items, nil
}

// StartArchive starts archiving the session for the given sessionID with the given
// name.
func (t *Tokbox) StartArchive(sessionID, name string) (Archive, error) {
	url := baseURL + "/" + apiVersion + "/project/" + t.key + "/archive/"

	var archive Archive
	if err := t.MakeRequest("POST", url, map[string]string{
		"sessionId":  sessionID,
		"name":       name,
		"outputMode": "composed",
	}, &archive); err != nil {
		return Archive{}, err
	}

	return archive, nil
}

// StopArchive stops the particualr archive for the given archiveID.
// any non-nil error denotes archiving stopped failed.
func (t *Tokbox) StopArchive(archiveID string) (Archive, error) {
	url := baseURL + "/" + apiVersion + "/project/" + t.key + "/archive/" + archiveID + "/stop/"
	var archive Archive
	if err := t.MakeRequest("POST", url, nil, &archive); err != nil {
		return Archive{}, err
	}
	return archive, nil
}

// NewRequest create a single http.Request based on url, headers and body that needs
// to be encoded. Returns the http.Response. any non-nil error means request is
// unsuccessfull.
func (t *Tokbox) MakeRequest(method, urlStr string, body interface{}, v interface{}) error {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, rel.String(), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	token, err := t.jwtToken()
	if err != nil {
		return err
	}
	req.Header.Set("X-OPENTOK-AUTH", token)

	res, err := http.DefaultClient.Do(req)

	if res.StatusCode >= 400 {
		return t.parseError(res)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// parseError returns the valid tokbox error if any of the response have status code
// greater than equal to 400.
func (t *Tokbox) parseError(resp *http.Response) error {
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	tErr := new(Error)
	if err = json.Unmarshal(resBody, tErr); err != nil {
		return errors.New("error in json body" + string(resBody))
	}
	return tErr
}
