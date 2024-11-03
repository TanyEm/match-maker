package lobby

import (
	"log"
	"sync"
	"time"

	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/TanyEm/match-maker/v2/internal/player"
)

const ErrNoMatch = "ErrNoMatch"

//go:generate mockgen -destination=./lobby_mock.go -package=lobby github.com/TanyEm/match-maker/v2/internal/lobby Lobbier
type Lobbier interface {
	AddPlayer(p player.Player)
	GetMatchByJoinID(joinID string) string
	GetMatchMakingTime() time.Duration
	Run()
	Stop()
}

type Lobby struct {
	stopCh          chan struct{}
	mu              sync.Mutex
	matchLocations  sync.Map
	WaitingTime     time.Duration
	MatchKeeper     match.Keeper
	playersToNotify map[string]string
}

func NewLobby(waitingTime time.Duration, matchKeeper match.Keeper) *Lobby {
	return &Lobby{
		stopCh:          make(chan struct{}),
		matchLocations:  sync.Map{},
		WaitingTime:     waitingTime,
		MatchKeeper:     matchKeeper,
		playersToNotify: make(map[string]string),
	}
}

func (l *Lobby) GetMatchMakingTime() time.Duration {
	return l.WaitingTime
}

func (l *Lobby) Run() {
	ticker := time.NewTicker(l.WaitingTime)
	defer ticker.Stop()

	log.Println("Lobby is running. Waiting for people to join...")

	for {
		select {
		case <-ticker.C:
			log.Println("Time is up! Start mathmaking...")
			l.StartMatches()
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
	matchLocation := &match.MatchLocation{}
	loaded, ok := l.matchLocations.Load(p.Country)
	if !ok {
		newMatch := match.NewMatch(p.Country, p.Level)
		newMatch.AddPlayer(p)

		matchLocation.Store(p.Level, newMatch)
		l.matchLocations.Store(p.Country, matchLocation)
		return
	}

	matchLocation = loaded.(*match.MatchLocation)

	// If the player's location is in the lobby, check if there is a match that the player can join.
	// I assume that the player can compete with players from the same level and the levels above and below
	// but the player with level 1 can only compete with players from levels 1, 2, and 3
	levelsToJoin := make([]int, 3)
	if p.Level > 1 {
		levelsToJoin[0], levelsToJoin[1], levelsToJoin[2] = p.Level-1, p.Level, p.Level+1
	} else {
		levelsToJoin[0], levelsToJoin[1], levelsToJoin[2] = 1, 2, 3
	}

	// If there is a match that the player can join, add the player to the match
	for _, level := range levelsToJoin {
		if loaded, ok := matchLocation.Load(level); ok {
			matchToJoin := loaded.(*match.Match)
			matchToJoin.AddPlayer(p)

			// If the match is full, start the match and delete it from the location in the lobby
			if matchToJoin.GetPlayersCount() == 10 {
				l.StartMatch(matchToJoin, matchLocation)
				matchLocation.Delete(level)
			}

			return
		}
	}

	// If no match is found, create a new match
	m := match.NewMatch(p.Country, p.Level)
	m.AddPlayer(p)

	// Store the match in the player's location
	matchLocation.Store(p.Level, m)
	l.matchLocations.Store(p.Country, matchLocation)
}

func (l *Lobby) StartMatch(m *match.Match, matchLocation *match.MatchLocation) {
	joinIDs := m.Start()
	l.mu.Lock()
	for _, joinID := range joinIDs {
		l.playersToNotify[joinID] = m.MatchID
	}
	l.mu.Unlock()

	leaderBoard := m.GetLeaderboard()
	l.MatchKeeper.AddLeaderBoard(&leaderBoard)
}

// StartMatches starts the matches that have more than one player across all locations in the lobby
// and cleans the lobby
func (l *Lobby) StartMatches() {
	l.matchLocations.Range(func(country, loaded interface{}) bool {
		matchLocation := loaded.(*match.MatchLocation)

		matchLocation.Range(func(level, loaded interface{}) bool {
			matchToStart := loaded.(*match.Match)

			// If there is more than one player in the match, start the match
			if matchToStart.GetPlayersCount() > 1 {
				l.StartMatch(matchToStart, matchLocation)
			} else {
				log.Printf("Match %s country %s level %d has only one player. Skipping the match and notifying the player...\n",
					matchToStart.MatchID,
					matchToStart.Country,
					matchToStart.Level,
				)

				l.mu.Lock()
				stalePlayer := matchToStart.GetPlayers()[0]
				l.playersToNotify[stalePlayer.JoinID] = ErrNoMatch
				l.mu.Unlock()
			}

			matchLocation.Delete(level)
			return true
		})

		l.matchLocations.Clear()
		return true
	})
}

func (l *Lobby) GetMatchByJoinID(joinID string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if matchID, ok := l.playersToNotify[joinID]; ok {
		return matchID
	}

	return ""
}
