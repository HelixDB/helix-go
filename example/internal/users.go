package internal

// Contains `User` related queries

import (
	"fmt"
	"log"

	"github.com/HelixDB/helix-go"
)

type User struct {
	ID        string
	Name      string
	Age       int32
	Email     string
	CreatedAt int32 `json:"created_at"`
	UpdatedAt int32 `json:"updated_at"`
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
		err := fmt.Errorf("Error while creating user: %s", err)
		return err
	}

	return nil
}

func CreateUsers(newUsers map[string]any) (map[string]any, error) {
	res, err := HelixClient.Query(
		"create_users",
		helix.WithData(newUsers),
	).AsMap()
	if err != nil {
		err = fmt.Errorf("Error while creating user: %s", err)
		return nil, err
	}

	return res, nil
}

func UpdateUser(userId string, newUserData map[string]any) error {
	_, err := HelixClient.Query(
		"update_user",
		helix.WithData(userId),
		helix.WithData(newUserData),
	).Raw()
	if err != nil {
		log.Printf("Error while creating user: %s", err)
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
		err = fmt.Errorf("Error while deleting user: %s", err)
		return err
	}

	return nil
}

func GetAllUsers(users *[]User) error {
	err := HelixClient.Query("get_users").Scan(
		helix.WithDest("users", &users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting users: %s", err)
		return err
	}

	return nil
}
