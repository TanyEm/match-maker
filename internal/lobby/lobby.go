package lobby

import (
	"log"
	"sync"
	"time"

	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/TanyEm/match-maker/v2/internal/player"
)

//go:generate mockgen -destination=./lobby_mock.go -package=lobby github.com/TanyEm/match-maker/v2/internal/lobby Lobbier
type Lobbier interface {
	AddPlayer(p player.Player)
	Run()
	Stop()
}

type Lobby struct {
	stopCh      chan struct{}
	mu          sync.Mutex
	players     sync.Map
	WaitingTime time.Duration
}

func NewLobby(waitingTime time.Duration) *Lobby {
	return &Lobby{
		stopCh:      make(chan struct{}),
		players:     sync.Map{},
		WaitingTime: waitingTime,
	}
}

func (l *Lobby) Run() {
	ticker := time.NewTicker(l.WaitingTime)
	defer ticker.Stop()

	log.Println("Lobby is running. Waiting for people to join...")

	for {
		select {
		case <-ticker.C:
			log.Println("Time is up! Cleaning up the lobby...")
			l.players.Clear()
		case <-time.After(100 * time.Millisecond): // Throttle the loop to decrease the load on the CPU
		case <-l.stopCh:
			log.Println("Lobby is stopped.")
			return
		}
	}
}

func (l *Lobby) Stop() {
	log.Println("Stopping the lobby...")
	l.stopCh <- struct{}{}
}

func (l *Lobby) AddPlayer(p player.Player) {
	log.Printf("Player %s joined the lobby, joinID: %s", p.PlayerID, p.JoinID)

	// If the player's location is not in the lobby, create a new match, new location and store it.
	matchLocation, ok := l.players.Load(p.Country)
	if !ok {
		m := match.NewMatch(p.Country, p.Level)
		m.AddPlayer(p)

		ml := &match.MatchLocation{}
		ml.Store(p.Level, m)

		l.players.Store(p.Country, ml)
		return
	}

	// If the player's location is in the lobby, check if there is a match that the player can join.
	// I assume that the player can compete with players from the same level and the levels above and below
	// but the player with level 1 can only compete with players from levels 1, 2, and 3
	matchLevelsToJoin := make([]int, 3)
	if p.Level > 1 {
		matchLevelsToJoin[0], matchLevelsToJoin[1], matchLevelsToJoin[2] = p.Level-1, p.Level, p.Level+1
	} else {
		matchLevelsToJoin[0], matchLevelsToJoin[1], matchLevelsToJoin[2] = 1, 2, 3
	}

	// If there is a match that the player can join, add the player to the match
	for _, level := range matchLevelsToJoin {
		if m, ok := matchLocation.(*match.MatchLocation).Load(level); ok {
			m.(*match.Match).AddPlayer(p)
			return
		}
	}

	// If no match is found, create a new match
	m := match.NewMatch(p.Country, p.Level)
	m.AddPlayer(p)

	// Store the match in the player's location
	matchLocation.(*match.MatchLocation).Store(p.Country, m)
	l.players.Store(p.Level, matchLocation)
}
