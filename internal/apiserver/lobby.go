package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LobbyRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	Level    int    `json:"level" binding:"min=1,max=99"`
	Country  string `json:"country" binding:"required,isocountry"`
}

type LobbyResponse struct {
	JoinID string `json:"join_id"`
}

func (s *APIServer) Lobby(ctx *gin.Context) {
	var req LobbyRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	joinID := uuid.New().String()
	ctx.JSON(http.StatusOK, LobbyResponse{JoinID: joinID})
}
