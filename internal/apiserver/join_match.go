package apiserver

import (
	"net/http"

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

	match := s.Match.Match(uuid.New().String())

	ctx.JSON(http.StatusOK, MatchResponse{MatchID: match.MatchID})
}
