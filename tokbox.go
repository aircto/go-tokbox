package tokbox

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
func (t *Tokbox) CreateSession() *Session {
	return nil
}

// jwtToken generates unique jwt token every time its called.
// Its used make any api request to tokbox REST API.
func (t *Tokbox) jwtToken() (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    t.key,
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Unix() + 180,
		Id:        string(uuid.NewV4()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secret))
}
