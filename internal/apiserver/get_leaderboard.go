package apiserver

import (
	"net/http"

	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetLeaderBoardResponse struct {
	match.LeaderBoard
}

func (s *APIServer) GetLeaderBoard(ctx *gin.Context) {
	matchID := ctx.Query("match_id")
	if matchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	if _, err := uuid.Parse(matchID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "match_id is not valid UUID"})
		return
	}

	leaderBoard := s.MatchKeeper.GetLeaderBoard(matchID)
	if leaderBoard == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "leaderboard not found"})
		return
	}

	leaderBoardResponse := GetLeaderBoardResponse{LeaderBoard: *leaderBoard}

	ctx.JSON(http.StatusOK, leaderBoardResponse)
}
