package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"go-docker/helpers"
	"io"
	"log"
)

func SendCodeToAContainer(code string, result chan string, Sender string) {

	if len(LiveContainers) == 0 {
		fmt.Println("No Live Containers")
		result <- "No Live Containers"
		return
	}

loop:
	for {
		for _, val := range LiveContainers {
			if *val.Workers > 0 {
				fmt.Println("val.Worker = ", *val.Workers, " At: ", val.Num)
				val.Receiver <- Job{
					Sender: Sender,
					Code:   code,
					result: result,
				}
				break loop
			} else {
				fmt.Println("Container Full: ", val.ID)
			}
		}
		//time.Sleep(1 * time.Second)
		fmt.Println("Looping at Run Code for: ", Sender)
	}
}

func ExecCode(code string, containerID string) string {

	exec, err := Client.Client.ContainerExecCreate(context.TODO(), containerID, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Tty:          false,
		Cmd:          []string{"python3", "main.py"},
	})

	execStartConfig := types.ExecStartCheck{}
	execReader, err := Client.Client.ContainerExecAttach(context.TODO(), exec.ID, execStartConfig)
	helpers.Check(err, "Running Exec Command")

	defer execReader.Close()

	var outBuffer bytes.Buffer
	_, err = io.Copy(&outBuffer, execReader.Reader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	body := outBuffer.String()
	helpers.Check(err, "Reading Exec Output to Body")

	return body
}

// docker ps -aq | xargs docker stop | xargs docker rm
