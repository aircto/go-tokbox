package tokbox

import "testing"

const (
	testKey    = "45822722"
	testSecret = "362f8bfbb5fff2f960c72ee4b798fa7029f9a601"
)

func TestTokBoxCreateSession(t *testing.T) {
	tok := New(testKey, testSecret)
	if tok.CreateSession() == nil {
		t.Errorf("expected non-nil, but got nil")
	}
}
