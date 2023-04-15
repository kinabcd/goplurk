package main

import (
	"fmt"
	"log"

	"github.com/kinabcd/goplurk"
)

func main() {
	var consumerToken = "..."
	var consumerSecret = "..."
	oauthRequest, err := goplurk.NewOAuthRequest(consumerToken, consumerSecret)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Open the following URL and authorize it:", oauthRequest.Url)

	var pinCode string
	fmt.Print("Input the PIN code: ")
	fmt.Scan(&pinCode)

	if client, token, secret, err := oauthRequest.SendPin(pinCode); err == nil {
		log.Println("AccessToken: " + token)
		log.Println("AccessSecret: " + secret)
		client.Users.Me()
	} else {
		log.Fatalf("Failed to create client: failed to SendPin: %v", err)
	}
}
