package main

import (
	"github.com/gin-gonic/gin"
	"go-docker/pkg/docker"
)

func addRoutes(eng *gin.Engine) {

	eng.GET("/liveCont", func(c *gin.Context) {

		ans := len(docker.LiveContainers)

		c.JSON(200, gin.H{
			"Count": ans,
		})

	})

}
