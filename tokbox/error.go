package tokbox

// Error represents any tokbox error return from REST API response.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error implements a error interface
func (e *Error) Error() string {
	return e.Message
}
