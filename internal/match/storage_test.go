package match

import (
	"sync"
	"testing"
)

func TestAddLeaderBoard(t *testing.T) {
	storage := NewStorage()
	lb := &LeaderBoard{MatchID: "match1"}

	storage.AddLeaderBoard(lb)

	if storage.GetLeaderBoard("match1") != lb {
		t.Errorf("Expected leaderboard to be added")
	}
}

func TestGetLeaderBoard(t *testing.T) {
	storage := NewStorage()
	lb := &LeaderBoard{MatchID: "match1"}

	storage.AddLeaderBoard(lb)

	retrievedLB := storage.GetLeaderBoard("match1")
	if retrievedLB != lb {
		t.Errorf("Expected to retrieve the correct leaderboard")
	}
}

func TestConcurrency(t *testing.T) {
	storage := NewStorage()
	lb := &LeaderBoard{MatchID: "match1"}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			storage.AddLeaderBoard(lb)
		}()
	}

	wg.Wait()

	retrievedLB := storage.GetLeaderBoard("match1")
	if retrievedLB != lb {
		t.Errorf("Expected to retrieve the correct leaderboard after concurrent writes")
	}
}
