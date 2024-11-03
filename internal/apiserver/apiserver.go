package apiserver

import (
	"net/http"

	"github.com/TanyEm/match-maker/v2/internal/lobby"
	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type APIServer struct {
	GinEngine *gin.Engine
	Lobby     lobby.Lobbier
	Match     match.Matcher
}

func NewAPIServer(lobby lobby.Lobbier) *APIServer {
	apiServer := &APIServer{
		Lobby: lobby,
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isocountry", ISOCountryValidator)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/lobby", apiServer.JoinLobby)
	r.GET("/match", apiServer.JoinMatch)

	apiServer.GinEngine = r

	return apiServer
}
