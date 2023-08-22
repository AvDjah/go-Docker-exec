package main

import (
	"fmt"
	"time"
)

type Job struct {
	Code   string
	Sender string
	Result chan string
}

func worker(ch chan Job, num int) {

	for {
		select {
		case job := <-ch:
			fmt.Println("Doing job: ", job.Sender, " at: ", num)
			time.Sleep(2 * time.Second)
			job.Result <- "Done"
		}
	}
}
