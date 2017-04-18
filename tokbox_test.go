package tokbox

import "testing"

const (
	testKey    = "45822722"
	testSecret = "362f8bfbb5fff2f960c72ee4b798fa7029f9a601"
)

func TestTokboxCreateSession(t *testing.T) {
	tok := New(testKey, testSecret)
	if tok.CreateSession() == nil {
		t.Errorf("expected non-nil, but got nil")
	}
}

func TestTokboxJwtToken(t *testing.T) {
	tok := New(testKey, testSecret)
	token1, err := tok.jwtToken()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	token2, err := tok.jwtToken()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if token1 == token2 {
		t.Errorf("expected unique token, got %v, %v", token1, token2)
	}
}
