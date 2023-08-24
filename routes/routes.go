package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-docker/helpers"
	"go-docker/pkg/docker"
	"net/http"
	"strconv"
	"sync/atomic"
)

func addRoutes(router *gin.Engine) {

	router.GET("/liveCont", func(c *gin.Context) {

		ans := len(docker.LiveContainers)

		c.JSON(200, gin.H{
			"Count": ans,
		})

	})

	router.PUT("/killContainer", func(c *gin.Context) {

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

	contNum := 1
	senders := int32(0)

	router.POST("/runCode", func(c *gin.Context) {
		atomic.AddInt32(&senders, 1)

		result := make(chan string, 2)

		docker.SendCodeToAContainer("Code", result, strconv.Itoa(int(senders)))

		fmt.Println("Waiting For Result at: ", senders)
		out := <-result

		fmt.Println(out)

		c.JSON(200, gin.H{
			"Result": out,
		})
	})

	router.POST("/createContainer", func(c *gin.Context) {

		list := client.ListContainers()

		if len(list) >= docker.ContainerLimit {
			c.JSON(http.StatusBadRequest, gin.H{
				"Result":      "Failure",
				"Description": "Container Limit Reacher. Consider Killing some off.",
			})
			return
		}

		client.StartContainer(&contNum)

		list = client.ListContainers()

		c.JSON(200, gin.H{
			"Result":     "Success",
			"Containers": list,
		})

	})

	router.GET("/listContainer", func(c *gin.Context) {

		list := client.ListContainers()

		c.JSON(200, gin.H{
			"Result": list,
		})

	})

	router.GET("/listImages", func(c *gin.Context) {
		list := client.ListImages()
		c.JSON(200, gin.H{
			"Result": list,
		})
	})

}
