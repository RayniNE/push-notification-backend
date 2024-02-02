package main

import (
	"fmt"
	"log"

	"github.com/raynine/push-notification/push_notification"
)

func main() {
	server := push_notification.NewServer()
	log.Println("Creating new server")

	server.GenerateVAPIDKeys()
	fmt.Println("Private Key: ", server.VAPIDPrivateKey)
	fmt.Println("Public Key: ", server.VAPIDPublicKey)

	server.Init()
	log.Println("Server running on port: 8080")
}
