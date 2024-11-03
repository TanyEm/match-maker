package match

import "sync"

//go:generate mockgen -destination=./storage_mock.go -package=match github.com/TanyEm/match-maker/v2/internal/match Keeper
type Keeper interface {
	AddLeaderBoard(lb *LeaderBoard)
	GetLeaderBoard(matchID string) *LeaderBoard
}

type Storage struct {
	matches map[string]*LeaderBoard
	mu      sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		matches: make(map[string]*LeaderBoard),
	}
}

func (s *Storage) AddLeaderBoard(lb *LeaderBoard) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches[lb.MatchID] = lb
}

func (s *Storage) GetLeaderBoard(matchID string) *LeaderBoard {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.matches[matchID]
}
