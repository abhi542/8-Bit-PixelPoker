package lobby

import (
	"testing"
)

func TestGenerateCode(t *testing.T) {
	c1 := GenerateCode()
	if len(c1) != 4 {
		t.Errorf("Expected length 4, got %d", len(c1))
	}
	c2 := GenerateCode()
	if c1 == c2 {
		t.Logf("Warning: Collision generated? %s == %s", c1, c2)
	}
}

func TestLobby_Join(t *testing.T) {
	cfg := LobbyConfig{MaxPlayers: 2}
	l := NewLobby("TEST", cfg)

	if err := l.Join("user1", "Alice"); err != nil {
		t.Errorf("First join failed: %v", err)
	}

	if err := l.Join("user2", "Bob"); err != nil {
		t.Errorf("Second join failed: %v", err)
	}

	if err := l.Join("user3", "Charlie"); err != ErrLobbyFull {
		t.Errorf("Expected ErrLobbyFull, got %v", err)
	}

	if l.State != StateReady {
		t.Errorf("Expected StateReady, got %v", l.State)
	}
}
