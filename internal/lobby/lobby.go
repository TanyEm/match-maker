package lobby

import (
	"encoding/json"
	"log"
	"sync"
	"time"

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
	players     []player.Player
	WaitingTime time.Duration
}

func NewLobby(waitingTime time.Duration) *Lobby {
	return &Lobby{
		stopCh:      make(chan struct{}),
		players:     make([]player.Player, 0),
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
			playersJSON, err := json.MarshalIndent(l.players, "", "  ")
			if err != nil {
				log.Printf("Error marshalling players: %v", err)
			} else {
				log.Println("Current players in the lobby:")
				log.Println(string(playersJSON))
			}

			log.Println("Time is up! Cleaning up the lobby...")
			l.mu.Lock()
			l.players = make([]player.Player, 0)
			l.mu.Unlock()
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
	l.mu.Lock()
	defer l.mu.Unlock()

	// Find the correct position to insert the player
	index := 0
	for i, existingPlayer := range l.players {
		if existingPlayer.Level < p.Level {
			index = i
			break
		}
		index = i + 1
	}

	// Insert the player at the correct position
	l.players = append(l.players[:index], append([]player.Player{p}, l.players[index:]...)...)
}
