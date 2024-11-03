package match

type LeaderBoard struct {
	MatchID string       `json:"match_id"`
	Players []PlayerInfo `json:"players"`
}

type PlayerInfo struct {
	PlayerID string `json:"player_id"`
	Level    int    `json:"level"`
	Country  string `json:"country"`
	Score    int    `json:"score"`
}
