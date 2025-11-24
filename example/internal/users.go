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
	res, err := HelixClient.Query(
		"create_user",
		helix.WithData(newUser),
	)
	if err != nil {
		return fmt.Errorf("Error while creating user: %v", err)
	}

	if err = res.Scan(user); err != nil {
		return fmt.Errorf("Error while scanning the \"create_user\" query result and populating `user`: %v", err)
	}

	return nil
}

func CreateUsers(newUsers map[string]any) error {
	_, err := HelixClient.Query(
		"create_users",
		helix.WithData(newUsers),
	)
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
	)
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
	)
	if err != nil {
		err = fmt.Errorf("Error while deleting user: %v", err)
		return err
	}

	return nil
}

// This function is meant to do 2 things
//   - populate the `users` slice
//   - return the helix response from the query
func GetAllUsers(users *[]User) (*helix.Response, error) {
	res, err := HelixClient.Query("get_users")
	if err != nil {
		return nil, fmt.Errorf("Error while getting users: %v", err)
	}

	if err = res.Scan(helix.WithDest("users", &users)); err != nil {
		return nil, fmt.Errorf("Error while scanning the \"users\" query result and populating the users slice: %v", err)
	}

	return res, nil
}
