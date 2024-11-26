package lobby

import (
	"testing"
	"time"

	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/TanyEm/match-maker/v2/internal/player"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewLobby(t *testing.T) {
	matchKeeper := match.NewMockKeeper(gomock.NewController(t))
	lobby := NewLobby(10*time.Second, matchKeeper)

	assert.NotNil(t, lobby)
	assert.Equal(t, 10*time.Second, lobby.GetMatchMakingTime())
	assert.NotNil(t, lobby.MatchKeeper)
	assert.NotNil(t, lobby.stopCh)
	isEmpty := true

	lobby.matchLocations.Range(func(_, _ interface{}) bool {
		isEmpty = false
		return false
	})

	assert.True(t, isEmpty, "Expected matchLocations to be empty")
	assert.NotNil(t, lobby.playersToNotify)
}

func TestLobby_AddSinglePlayer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMatchKeeper := match.NewMockKeeper(mockCtrl)
	l := NewLobby(1*time.Minute, mockMatchKeeper)

	tests := []struct {
		name            string
		player          player.Player
		expectedCountry string
		expectedLevel   int
	}{
		{
			name: "New player with level 1",
			player: player.Player{
				PlayerID: "1",
				JoinID:   "3525d198-e0a3-40c1-9689-c67d75ea6f63",
				Level:    1,
				Country:  "FIN",
			},
			expectedCountry: "FIN",
			expectedLevel:   2, // Level is set to 2 for level 1 players
		},
		{
			name: "New player from different country",
			player: player.Player{
				PlayerID: "3",
				JoinID:   "",
				Level:    2,
				Country:  "USA",
			},
			expectedCountry: "USA",
			expectedLevel:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l.AddPlayer(tt.player)

			location, ok := l.matchLocations.Load(tt.expectedCountry)
			assert.True(t, ok, "Expected match location for country %s", tt.expectedCountry)

			matchLocation := location.(*match.MatchLocation)
			loadedMatch, ok := matchLocation.Load(tt.expectedLevel)
			assert.True(t, ok, "Expected match at level %d", tt.expectedLevel)

			m := loadedMatch.(*match.Match)
			assert.Equal(t, 1, m.GetPlayersCount(), "Expected player count to be 1")
			assert.Equal(t, tt.player.PlayerID, m.GetPlayers()[0].PlayerID, "Expected player ID to match")
		})
	}
}

func TestLobby_AddMultiplePlayers2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockKeeper := match.NewMockKeeper(mockCtrl)

	mockKeeper.EXPECT().
		AddLeaderBoard(gomock.Any()).
		Times(1)

	l := NewLobby(1*time.Minute, mockKeeper)

	tests := []struct {
		name                       string
		players                    []player.Player
		expectedPlayerCountInLobby int
		expectedMatchCount         int
	}{
		{
			name: "Add 11 players to match, with only 10 fitting",
			players: []player.Player{
				{PlayerID: "player1", JoinID: "join10", Country: "FIN", Level: 1},
				{PlayerID: "player2", JoinID: "join20", Country: "FIN", Level: 2},
				{PlayerID: "player3", JoinID: "join30", Country: "FIN", Level: 3},
				{PlayerID: "player4", JoinID: "join40", Country: "FIN", Level: 1},
				{PlayerID: "player5", JoinID: "join1", Country: "FIN", Level: 1},
				{PlayerID: "player6", JoinID: "join2", Country: "FIN", Level: 2},
				{PlayerID: "player7", JoinID: "join3", Country: "FIN", Level: 3},
				{PlayerID: "player8", JoinID: "join4", Country: "FIN", Level: 1},
				{PlayerID: "player9", JoinID: "join5", Country: "FIN", Level: 1},
				{PlayerID: "player10", JoinID: "join10", Country: "FIN", Level: 3},
				{PlayerID: "player11", JoinID: "join11", Country: "FIN", Level: 1},
			},
			expectedPlayerCountInLobby: 1, // Only 10 players fit in the match, so the 11th player is left in the lobby
			expectedMatchCount:         1, // All players are from the same country "FIN"
		},
		{
			name: "New player from different country",
			players: []player.Player{
				{PlayerID: "player1", JoinID: "join1", Country: "FIN", Level: 1},
				{PlayerID: "player22", JoinID: "join2", Country: "USA", Level: 1},
			},
			expectedPlayerCountInLobby: 2, // Both players should be matched
			expectedMatchCount:         2, // One match per country
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = NewLobby(1*time.Minute, mockKeeper)

			for _, p := range tt.players {
				l.AddPlayer(p)
			}

			// Count the total number of players matched and the number of distinct matches
			totalPlayersCount := 0
			totalMatchesCount := 0
			l.matchLocations.Range(func(_, location interface{}) bool {
				totalMatchesCount++
				matchLocation := location.(*match.MatchLocation)
				matchLocation.Range(func(_, m interface{}) bool {
					match := m.(*match.Match)
					totalPlayersCount += match.GetPlayersCount()
					return true
				})
				return true
			})

			assert.Equal(t, tt.expectedPlayerCountInLobby, totalPlayersCount, "Expected player count to match")
			assert.Equal(t, tt.expectedMatchCount, totalMatchesCount, "Expected match count to match")
		})
	}
}

func TestLobby_GetMatchByJoinID(t *testing.T) {
	matchKeeper := match.NewMockKeeper(gomock.NewController(t))
	lobby := NewLobby(10*time.Second, matchKeeper)

	player1 := player.Player{PlayerID: "player1", JoinID: "join1", Country: "FIN", Level: 1}
	lobby.AddPlayer(player1)

	lobby.mu.Lock()
	lobby.playersToNotify["join1"] = "match1"
	lobby.mu.Unlock()

	matchID := lobby.GetMatchByJoinID("join1")
	assert.Equal(t, "match1", matchID)

	matchID = lobby.GetMatchByJoinID("nonexistent")
	assert.Equal(t, "", matchID)
}

func TestLobby_GetMatchByJoinID2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMatchKeeper := match.NewMockKeeper(mockCtrl)
	l := NewLobby(1*time.Minute, mockMatchKeeper)

	// Add players to the lobby and simulate match creation
	player1 := player.Player{PlayerID: "1", JoinID: "join1", Level: 2, Country: "FIN"}
	player2 := player.Player{PlayerID: "2", JoinID: "join2", Level: 3, Country: "FIN"}

	l.AddPlayer(player1)
	l.AddPlayer(player2)

	matchID := "match123"
	l.playersToNotify["join1"] = matchID

	tests := []struct {
		name     string
		joinID   string
		expected string
	}{
		{"Valid joinID", "join1", matchID},
		{"Invalid joinID", "join3", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := l.GetMatchByJoinID(tt.joinID)
			assert.Equal(t, tt.expected, result, "Expected match ID to match")
		})
	}
}
