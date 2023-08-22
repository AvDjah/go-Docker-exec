package helpers

import (
	"fmt"
	"log"
)

func Check(err error, msg string) {
	if err != nil {
		log.Panic("Error at: ", msg, " : ", err)
		return
	} else {
		fmt.Println("No Err: ", msg)
	}
}
