package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-docker/helpers"
	"go-docker/pkg/docker"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var ContainerLimit = 4

var ContainerID = "5202db9b554dbac6148099e11c389fa6863cb9ee0b9653da0775a4a3394aa8d0"

func main() {
	client := docker.New()
	defer func() {
		err := client.Client.Close()
		helpers.Check(err, "Closing Client")
	}()

	client.ListContainers()
	client.Init()
	contNum := 1

	router := gin.Default()

	ch := make(chan Job, 20)

	for i := 1; i < 20; i++ {
		go worker(ch, i)
	}

	i := 0
	served := 0

	senders := int32(0)

	addRoutes(router)

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

		if len(list) >= ContainerLimit {
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

	router.POST("/runFiles", func(c *gin.Context) {

		mp := gin.H{}

		err := c.Bind(&mp)
		helpers.Check(err, "Extracting Body to Map")

		for i := range mp {
			fmt.Println(i, mp[i])
		}

		myCode := i
		i += 1

		receiver := make(chan string)

	loop:
		for {
			select {

			case ch <- Job{Sender: strconv.Itoa(myCode),
				Code:   "Code",
				Result: receiver}:
				break loop

			default:
				fmt.Println("Channel Full Waiting at: ", myCode)
				time.Sleep(1 * time.Second)
			}
		}

		var out string
		result := <-receiver
		out = result

		//singleWorker()
		served += 1
		fmt.Println("Served: ", served)
		c.JSON(200, gin.H{
			"Result": out,
		})
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
		return
	}
}
