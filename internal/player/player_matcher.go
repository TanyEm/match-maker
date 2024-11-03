package player

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// PlayerMatcher is a custom matcher for player.Player
type PlayerMatcher struct {
	expected Player
}

// Matches checks if the actual player matches the expected player
func (m PlayerMatcher) Matches(x interface{}) bool {
	actual, ok := x.(Player)
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
func EqPlayer(p Player) gomock.Matcher {
	return PlayerMatcher{expected: p}
}
