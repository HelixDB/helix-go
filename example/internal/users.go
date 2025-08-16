package internal

// Contains `User` related queries

import (
	"fmt"
	"github.com/HelixDB/helix-go"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int32  `json:"age"`
	Email     string `json:"email"`
	CreatedAt int32  `json:"created_at"`
}

// Create a type struct for the "create_user" query
type CreateUserResponse struct {
	User User `json:"user"`
}

func CreateUser(newUser map[string]any, user *CreateUserResponse) error {
	err := HelixClient.Query(
		"create_user",
		helix.WithData(newUser),
	).Scan(user)
	if err != nil {
		err := fmt.Errorf("Error while creating user: %v", err)
		return err
	}

	return nil
}

func CreateUsers(newUsers map[string]any) error {
	_, err := HelixClient.Query(
		"create_users",
		helix.WithData(newUsers),
	).Raw()
	if err != nil {
		err = fmt.Errorf("Error while creating users: %v", err)
		return err
	}

	return nil
}

func UpdateUser(user User) error {
	_, err := HelixClient.Query(
		"update_user",
		helix.WithData(user),
	).Raw()
	if err != nil {
		err = fmt.Errorf("Error while updating user: %s", err)
		return err
	}

	return nil
}

func DeleteUser(data map[string]any) error {
	_, err := HelixClient.Query(
		"delete_user",
		helix.WithData(data),
	).Raw()
	if err != nil {
		err = fmt.Errorf("Error while deleting user: %v", err)
		return err
	}

	return nil
}

func GetAllUsers(users *[]User) error {
	err := HelixClient.Query("get_users").Scan(
		helix.WithDest("users", &users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting users: %v", err)
		return err
	}

	return nil
}
