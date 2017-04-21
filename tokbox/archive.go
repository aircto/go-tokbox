package tokbox

// Archive denotes single recorded archive of any session.
type Archive struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"` // unix timestamp
	Duration  int64  `json:"duration"`   // in secs
	HasAudio  bool   `json:"has_audio"`
	HasVideo  bool   `json:"has_video"`
	Name      string `json:"name"`
	ProjectID int64  `json:"projectID"`
	Reason    string `json:"reason"`
	SessionID string `json:"session_id"`
	Size      int64  `json:"size"`
	Status    string `json:"status"`
	URL       string `json:"url"`
}

// ArchiveList represents single response of archive list endpoint
type ArchiveList struct {
	Count int       `json:"count"`
	Items []Archive `json:"items"`
}
