package main

import (
	"fmt"
	"log"
	"time"

	helix "github.com/HelixDB/helix-go"
)

var HelixClient *helix.Client

func main() {

	// Connect to client
	HelixClient = helix.NewClient("http://localhost:6969")

	// Create user data

	now := time.Now()

	timestamp := now.Unix()

	timestamp32 := int32(timestamp)

	newUser := map[string]any{
		"name":  "John",
		"age":   21,
		"email": "johndoe@email.com",
		"now":   timestamp32,
	}

	// Create user in Helix
	createUserRes, err := HelixClient.Query(
		"create_user",
		helix.WithData(newUser),
	)
	if err != nil {
		log.Fatalf("Error while creating user: %s", err)
	}

	fmt.Println(createUserRes)

	// Get all users
	users, err := HelixClient.Query("get_users")
	if err != nil {
		log.Fatalf("Error while creating user: %s", err)
	}

	fmt.Println(users)
}
