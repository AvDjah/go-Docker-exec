package routes

import (
	"github.com/gin-gonic/gin"
	"go-docker/pkg/docker"
)

var router *gin.Engine
var client *docker.ClientType

func New() *gin.Engine {

	if router != nil {
		return router
	}

	if client == nil {
		client = docker.New()
	}

	eng := gin.Default()
	addRoutes(eng)

	return eng
}
