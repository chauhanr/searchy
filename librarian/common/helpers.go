package common

import (
	"log"
	"fmt"
)

func Log(msg string){
	log.Println("INFO: ", msg)
}

func Warn(msg string){
	log.Println("---------------------------------")
	log.Println(fmt.Sprintf("WARN: %s\n", msg))
	log.Println("---------------------------------")
}
