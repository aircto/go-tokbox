package tokbox

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// Session represents a single session for the clients to communicate.
type Session struct {
	ID             string    `json:"session_id"`
	ProjectID      string    `json:"project_id"`
	CreatedAt      time.Time `json:"created_dt"` // as per REST API Doc. https://tokbox.com/developer/rest/#session_id_production
	MediaServerURL string    `json:"media_server_url"`
}

func (s *Session) Token(key, secret string) (string, error) {
	now := time.Now().Unix()
	expire := now + 3600
	ran := rand.Intn(999999)

	payload := fmt.Sprintf("create_time=%d&session_id=%s&nonce=%d&expire_time=%d&connection_data=None&role=publisher", now, s.ID, ran, expire)

	sig := s.sign(payload, secret)

	data := fmt.Sprintf("partner_id=%s&sig=%s:%s", key, sig, payload)

	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return "T1==" + encoded, nil
}

func (s *Session) sign(payload, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(payload))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
