package main

import (
	"go-docker/helpers"
	"go-docker/pkg/docker"
	"go-docker/routes"
	"log"
	"net/http"
)

func main() {
	client := docker.New()
	defer func() {
		err := client.Client.Close()
		helpers.Check(err, "Closing Client")
	}()

	client.ListContainers()
	client.Init()

	router := routes.New()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
		return
	}
}
