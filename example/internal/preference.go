package internal

import (
	"fmt"

	"github.com/HelixDB/helix-go"
)

func CreatePreference(data map[string]any) error {
	_, err := HelixClient.Query(
		"create_preference",
		helix.WithData(data),
	)
	if err != nil {
		err := fmt.Errorf("Error while creating preference: %v", err)
		return err
	}

	return nil
}

func AddPreferenceToUser(data map[string]any) error {
	_, err := HelixClient.Query(
		"add_preference_to_user",
		helix.WithData(data),
	)
	if err != nil {
		err := fmt.Errorf("Error while adding preference to user: %v", err)
		return err
	}

	return nil
}

// This function is meant to do 2 things
//   - populate the `users` slice
//   - return the body in form of bytes ([]byte) from the query
func SearchUsersByPreference(data map[string]any, users *[]User) ([]byte, error) {
	res, err := HelixClient.Query(
		"search_users_by_preference",
		helix.WithData(data),
	)
	if err != nil {
		return nil, fmt.Errorf("Error while searching for users by preferences: %v", err)
	}

	if err = res.Scan(helix.WithDest("users", users)); err != nil {
		return nil, fmt.Errorf("Error while scanning the \"search_users_by_preference\" query result and populating the users slice: %v", err)
	}

	body := res.Raw()

	return body, nil
}
