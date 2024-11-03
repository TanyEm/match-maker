package match

type LeaderBoard struct {
	MatchID string
	Players []PlayerInfo
}

type PlayerInfo struct {
	PlayerID string
	Level    int
	Country  string
	Score    int
}
