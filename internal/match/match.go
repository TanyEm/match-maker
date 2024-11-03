package match

import (
	"log"
	"sync"

	"github.com/TanyEm/match-maker/v2/internal/player"
	"github.com/google/uuid"
)

// MatchLocation a map of country to a list of matches
type MatchLocation struct {
	sync.Map
}

type Matcher interface {
	Match(matchID string) *Match
	AddPlayer(p player.Player)
	GetPlayersCount() int
	Start()
}

type Match struct {
	MatchID string
	Level   int
	Country string
	players []player.Player
	mu      sync.Mutex
	started bool
}

func NewMatch(country string, level int) *Match {
	return &Match{
		MatchID: uuid.New().String(),
		Level:   level,
		Country: country,
		players: []player.Player{},
		started: false,
	}
}

func (m *Match) AddPlayer(p player.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the correct position to insert the player
	index := 0
	for i, existingPlayer := range m.players {
		if existingPlayer.Level < p.Level {
			index = i
			break
		}
		index = i + 1
	}

	// Insert the player at the correct position
	m.players = append(m.players[:index], append([]player.Player{p}, m.players[index:]...)...)
	log.Printf("Player %s level %d joined the match with %d people: country %s level %d matchID: %s", p.PlayerID, p.Level, len(m.players), m.Country, m.Level, m.MatchID)
}

func (m *Match) Start() []string {
	m.started = true
	log.Printf("Match %s started. Notifying %d players...", m.MatchID, len(m.players))

	joinIDs := make([]string, 0, len(m.players))

	for _, p := range m.players {
		joinIDs = append(joinIDs, p.JoinID)
	}

	return joinIDs
}

func (m *Match) GetPlayersCount() int {
	return len(m.players)
}

func (m *Match) GetLeaderboard() LeaderBoard {
	leaderBoard := LeaderBoard{
		MatchID: m.MatchID,
		Players: make([]PlayerInfo, 0, len(m.players)),
	}

	for _, p := range m.players {
		playerInfo := PlayerInfo{
			PlayerID: p.PlayerID,
			Level:    p.Level,
			Country:  p.Country,
			Score:    0,
		}
		leaderBoard.Players = append(leaderBoard.Players, playerInfo)
	}

	return leaderBoard
}
