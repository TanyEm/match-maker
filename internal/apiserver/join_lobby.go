package apiserver

import (
	"net/http"

	"github.com/TanyEm/match-maker/v2/internal/player"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LobbyRequest is a request to join a lobby
// I assume that
// - playerID is any string
// - players level is a number between 1 and 99
// - country is a valid ISO 3 letters country code
type LobbyRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	Level    int    `json:"level" binding:"min=1,max=99"`
	Country  string `json:"country" binding:"required,isocountry"`
}

type LobbyResponse struct {
	JoinID string `json:"join_id"`
}

func (s *APIServer) JoinLobby(ctx *gin.Context) {
	var req LobbyRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := player.Player{
		PlayerID: req.PlayerID,
		Level:    req.Level,
		Country:  req.Country,
		JoinID:   uuid.New().String(),
	}

	s.Lobby.AddPlayer(player)

	ctx.JSON(http.StatusOK, LobbyResponse{JoinID: player.JoinID})
}
