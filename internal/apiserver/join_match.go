package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchResponse struct {
	MatchID string `json:"match_id"`
}

func (s *APIServer) JoinMatch(ctx *gin.Context) {
	// Create a context with a 30-second timeout
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Replace the request's context with the new one
	ctx.Request = ctx.Request.WithContext(c)

	joinID := ctx.Query("join_id")
	if joinID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "join_id is required"})
		return
	}

	if _, err := uuid.Parse(joinID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "join_id is not valid UUID"})
		return
	}

	var matchID string
	for {
		matchID = s.Lobby.GetMatchByJoinID(joinID)
		if matchID != "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	ctx.JSON(http.StatusOK, MatchResponse{MatchID: matchID})
}
