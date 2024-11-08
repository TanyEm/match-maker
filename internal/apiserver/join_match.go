package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/TanyEm/match-maker/v2/internal/lobby"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchResponse struct {
	MatchID string `json:"match_id"`
}

func (s *APIServer) JoinMatch(ctx *gin.Context) {
	joinID := ctx.Query("join_id")
	if joinID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "join_id is required"})
		return
	}

	if _, err := uuid.Parse(joinID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "join_id is not valid UUID"})
		return
	}

	// Create a context with a s.Lobby.GetMatchMakingTime() (default 30 sec) timeout
	c, cancel := context.WithTimeout(ctx.Request.Context(), s.Lobby.GetMatchMakingTime())
	defer cancel()

	// Replace the request's context with the new one
	ctx.Request = ctx.Request.WithContext(c)

	var matchID string

	// Wait for the match to be created for s.Lobby.GetMatchMakingTime() seconds
	for {
		matchID = s.Lobby.GetMatchByJoinID(joinID)
		if matchID != "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
		if c.Err() != nil {
			matchID = lobby.ErrNoMatch
			break
		}
	}

	if matchID == lobby.ErrNoMatch {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no match for the player, try to join the lobby again"})
		return
	}

	ctx.JSON(http.StatusOK, MatchResponse{MatchID: matchID})
}
