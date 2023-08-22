package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	client2 "github.com/docker/docker/client"
	"go-docker/helpers"
	"strconv"
	"sync"
	"sync/atomic"
)

type ClientType struct {
	Client        *client2.Client
	ContainerName string
}

type Job struct {
	Sender string
	Code   string
	result chan string
}

type Container struct {
	Num      int
	ID       string
	Receiver chan Job
	Workers  *int32
	mu       sync.Mutex
	Kill     chan string
}

var ImageName = "leetimg"

var LiveContainers []*Container

var Client *ClientType

func New() *ClientType {

	if Client != nil {
		fmt.Println("Already Have Client")
		return Client
	}

	NewClient, err := client2.NewClientWithOpts(client2.FromEnv, client2.WithAPIVersionNegotiation())
	Client = &ClientType{
		Client:        NewClient,
		ContainerName: "Container",
	}
	helpers.Check(err, "Creating Docker Client")
	return Client
}

func (cli *ClientType) ListContainers() []string {

	containerList, err := cli.Client.ContainerList(context.TODO(), types.ContainerListOptions{})
	helpers.Check(err, "Getting Container List")

	var arr []string

	for _, val := range containerList {
		fmt.Println(val.Names[0], " : ", val.ID)
		arr = append(arr, val.Names[0])
	}

	return arr
}

func (cli *ClientType) ListImages() []string {
	imageList, err := cli.Client.ImageList(context.TODO(), types.ImageListOptions{})
	helpers.Check(err, "Getting Image List")

	var arr []string

	for _, val := range imageList {
		fmt.Println(val.RepoTags[0])
		arr = append(arr, val.RepoTags[0])
	}

	return arr
}

func (cli *ClientType) StartContainer(num *int) {

	ctx := context.Background()
	var contID string

	for {
		contName := cli.ContainerName + strconv.Itoa(*num)
		res, err := cli.Client.ContainerCreate(ctx, &container.Config{
			Image: ImageName,
			Tty:   true,
		}, nil, nil, nil, contName)
		if err == nil {
			contID = res.ID
			fmt.Println("Started Container: ", contName)
			go createContainerWorkers(res.ID, *num)
			break
		} else {
			fmt.Println("ISSUE RUNNING CONTAINER FOR N:", *num)
			*num += 1
		}
	}

	err := cli.Client.ContainerStart(ctx, contID, types.ContainerStartOptions{})
	helpers.Check(err, "Starting Container")
}

func createContainerWorkers(ID string, num int) {

	receiver := make(chan Job, 10)

	killContainer := make(chan string)

	workers := int32(5)

	newContainer := Container{
		Num:      num,
		ID:       ID,
		Receiver: receiver,
		Workers:  &workers,
		Kill:     killContainer,
	}

	for i := 1; i <= 5; i++ {
		go newContainer.worker(i, receiver, num)
	}

	LiveContainers = append(LiveContainers, &newContainer)

	<-killContainer
}

func (cont *Container) worker(num int, receiver chan Job, contNum int) {

	for {
		select {
		case job := <-receiver:

			atomic.AddInt32(cont.Workers, -1)

			out := ExecCode(job.Code, cont.ID)
			job.result <- "Done: " + out
			//job.result <- "Job Done from worker: " + strconv.Itoa(num) + " at: " + strconv.Itoa(contNum)
			atomic.AddInt32(cont.Workers, 1)
		}
	}

}
