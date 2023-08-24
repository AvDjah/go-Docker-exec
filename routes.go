package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-docker/helpers"
	"go-docker/pkg/docker"
)

func addRoutes(eng *gin.Engine) {

	eng.GET("/liveCont", func(c *gin.Context) {

		ans := len(docker.LiveContainers)

		c.JSON(200, gin.H{
			"Count": ans,
		})

	})

	eng.PUT("/killContainer", func(c *gin.Context) {

		body := make(map[string]interface{})
		err := c.BindJSON(&body)
		helpers.Check(err, "Binding JSON")

		containerName := fmt.Sprint(body["Name"])

		client := docker.New()
		go client.KillContainer("/" + containerName)

		c.JSON(200, gin.H{
			"Body": containerName,
		})
	})

}
