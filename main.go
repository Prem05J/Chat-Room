package main

import (
	"log"

	application "github.com/chatroom-go/Application"
)

func main() {
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
