package tokbox

import "time"

// Session represents a single session for the clients to communicate.
type Session struct {
	ID             string    `json:"session_id"`
	ProjectID      string    `json:"project_id"`
	CreatedAt      time.Time `json:"created_dt"` // as per REST API Doc. https://tokbox.com/developer/rest/#session_id_production
	MediaServerURL string    `json:"media_server_url"`
}
