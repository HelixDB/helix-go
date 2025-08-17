package main

import (
	"fmt"
	"log"
	"time"

	"example/internal"
)

func main() {
	// Initialize Helix client
	internal.ConfigHelix()
	fmt.Println("✓ Helix client initialized")

	now := time.Now()
	timestamp32 := int32(now.Unix())

	// Create a user
	fmt.Println("\n--- Creating first user ---")
	newUser := map[string]any{
		"name":  "John Doe",
		"age":   25,
		"email": "johndoe@email.com",
		"now":   timestamp32,
	}

	var createUserResponse internal.CreateUserResponse

	err := internal.CreateUser(newUser, &createUserResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create user response: %+v\n", createUserResponse)

	// Create 2 more users
	fmt.Println("\n--- Creating 2 more users ---")
	user2 := map[string]any{
		"name":  "Jane Smith",
		"age":   28,
		"email": "janesmith@email.com",
		"now":   timestamp32,
	}
	user3 := map[string]any{
		"name":  "Bob Wilson",
		"age":   32,
		"email": "bobwilson@email.com",
		"now":   timestamp32,
	}

	err = internal.CreateUsers(
		map[string]any{
			"users": []map[string]any{user2, user3},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success creating users\n")

	// Get all users
	fmt.Println("\n--- Retrieving all users ---")
	var users []internal.User
	err = internal.GetAllUsers(&users)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Users:\n")
	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

	// Add follow relationships
	fmt.Println("\n--- Creating follow relationships ---")

	followInput1 := &internal.FollowUserInput{
		FollowerId: users[0].ID,
		FollowedId: users[1].ID,
	}
	err = internal.FollowUser(followInput1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s follows %s\n", users[0].Name, users[1].Name)

	followInput2 := &internal.FollowUserInput{
		FollowerId: users[1].ID,
		FollowedId: users[2].ID,
	}
	err = internal.FollowUser(followInput2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s follows %s\n", users[1].Name, users[2].Name)

	followInput3 := &internal.FollowUserInput{
		FollowerId: users[2].ID,
		FollowedId: users[0].ID,
	}
	err = internal.FollowUser(followInput3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s follows %s\n", users[2].Name, users[0].Name)

	followInput4 := &internal.FollowUserInput{
		FollowerId: users[0].ID,
		FollowedId: users[2].ID,
	}
	err = internal.FollowUser(followInput4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s follows %s\n", users[0].Name, users[2].Name)

	followInput5 := &internal.FollowUserInput{
		FollowerId: users[1].ID,
		FollowedId: users[0].ID,
	}
	err = internal.FollowUser(followInput5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s follows %s\n", users[1].Name, users[0].Name)

	fmt.Println("\n--- User followers and following ---")
	for _, user := range users {
		fmt.Printf("\nUser: %s\n", user.Name)
		var followers []internal.User
		err := internal.Followers(
			map[string]any{
				"id": user.ID,
			},
			&followers,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\tFollowers:")

		for _, follower := range followers {
			fmt.Printf("\t\t%s\n", follower.Name)
		}

		var following []internal.User
		err = internal.Following(
			map[string]any{
				"id": user.ID,
			},
			&following,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\tFollowing:")

		for _, userFollowing := range following {
			fmt.Printf("\t\t%s\n", userFollowing.Name)
		}
	}

	// Update the last user from `users`
	userToUpdate := users[len(users)-1]
	fmt.Printf("\n--- Update user: %s ---\n", userToUpdate.Name)
	userToUpdate.Age = 23
	err = internal.UpdateUser(userToUpdate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User %s successfully updated\n", userToUpdate.Name)
	fmt.Printf("Old age: %d; New age: %d\n", users[len(users)-1].Age, userToUpdate.Age)

	users[len(users)-1].Age = userToUpdate.Age

	// Create preference
	newPreference := "cats"
	newPreference2 := "dogs"
	fmt.Printf("\n--- Create preferences \"%s\" and \"%s\" ---\n", newPreference, newPreference2)
	err = internal.CreatePreference(
		map[string]any{"preference": newPreference},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New preference \"%s\" created successfully\n", newPreference)

	err = internal.CreatePreference(
		map[string]any{"preference": newPreference2},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New preference \"%s\" created successfully\n", newPreference2)

	// Add preference to user
	fmt.Printf("\n--- Add preference \"cats\" to user: %s and %s ---\n", users[0].Name, users[1].Name)
	err = internal.AddPreferenceToUser(
		map[string]any{
			"preference": "cats",
			"user_id":    users[0].ID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Preference \"cats\" successfully added to user %s\n", users[0].Name)

	err = internal.AddPreferenceToUser(
		map[string]any{
			"preference": "cats",
			"user_id":    users[1].ID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Preference \"cats\" successfully added to user %s\n", users[1].Name)

	err = internal.AddPreferenceToUser(
		map[string]any{
			"preference": "dogs",
			"user_id":    users[2].ID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Preference \"dogs\" successfully added to user %s\n", users[2].Name)

	// Search users by preference
	fmt.Println("\n--- Search Users by Preference: \"cats\" ---")
	var usersSearchResults []internal.User
	err = internal.SearchUsersByPreference(
		map[string]any{
			"preference": "cats",
			"limit":      5,
		},
		&usersSearchResults,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Users Search Result:\n")
	for _, user := range usersSearchResults {
		fmt.Printf("%+v\n", user.Name)
	}

	fmt.Println("Users search by preference completed successfully")

	// Delete the first user from `users`
	fmt.Printf("\n--- Delete User: %s ---", users[0].Name)
	err = internal.DeleteUser(
		map[string]any{"id": users[0].ID},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nUser successfully deleted")

	fmt.Println("\nExample completed successfully!")
}
