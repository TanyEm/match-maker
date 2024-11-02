package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type APIServer struct {
	GinEngine *gin.Engine
}

func NewAPIServer() *APIServer {
	apiServer := &APIServer{}

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
	r.POST("/lobby", apiServer.Lobby)
	r.GET("/match", apiServer.Match)

	apiServer.GinEngine = r

	return apiServer
}
