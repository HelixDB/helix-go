package main

import (
	"fmt"
	"log"
	"time"

	helix "github.com/HelixDB/helix-go"
)

var HelixClient *helix.Client

// Create user struct
type User struct {
	Name      string
	Age       int32
	Email     string
	CreatedAt int32 `json:"created_at"`
	UpdatedAt int32 `json:"updated_at"`
}

// Create a type struct for the "get_users" query
type GetUsersResponse struct {
	Users []User `json:"users"`
}

// Create a type struct for the "create_users" query
type CreateUserResponse struct {
	User []User `json:"user"`
}

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
	var createdUser CreateUserResponse
	err := HelixClient.Query(
		"create_user",
		helix.WithData(newUser),
	).Scan(&createdUser)
	if err != nil {
		log.Fatalf("Error while creating user: %s", err)
	}

	fmt.Println(createdUser)

	// Get all users and put Helix's response in GetUsersResponse
	var getUsersResponse GetUsersResponse
	err = HelixClient.Query("get_users").Scan(&getUsersResponse)
	if err != nil {
		log.Fatalf("Error while getting users: %s", err)
	}

	fmt.Println(getUsersResponse)

	// Get all users and put "users" from Helix's response in the `users` variable
	var users []User
	err = HelixClient.Query("get_users").Scan(
		helix.WithDest("users", &users),
	)
	if err != nil {
		log.Fatalf("Error while getting users: %s", err)
	}

	fmt.Println(users)

	// Get all users in go's `map` data type
	usersMap, err := HelixClient.Query("get_users").AsMap()
	if err != nil {
		log.Fatalf("Error while getting users: %s", err)
	}

	fmt.Println(usersMap["users"])
}
