package utils

import (
	"fmt"

	"github.com/TanyEm/match-maker/v2/internal/player"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// PlayerMatcher is a custom matcher for player.Player
type PlayerMatcher struct {
	expected player.Player
}

// Matches checks if the actual player matches the expected player
func (m PlayerMatcher) Matches(x interface{}) bool {
	actual, ok := x.(player.Player)
	if !ok {
		return false
	}

	if _, err := uuid.Parse(actual.JoinID); err != nil {
		return false
	}

	return m.expected.PlayerID == actual.PlayerID &&
		m.expected.Level == actual.Level &&
		m.expected.Country == actual.Country
}

// String returns a string representation of the matcher
func (m PlayerMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

// EqPlayer returns a PlayerMatcher for the given player
func EqPlayer(p player.Player) gomock.Matcher {
	return PlayerMatcher{expected: p}
}
