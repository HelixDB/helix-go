package internal

import (
	"fmt"

	"github.com/HelixDB/helix-go"
)

func CreatePreference(data map[string]any) error {
	_, err := HelixClient.Query(
		"create_preference",
		helix.WithData(data),
	).Raw()
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
	).Raw()
	if err != nil {
		err := fmt.Errorf("Error while adding preference to user: %v", err)
		return err
	}

	return nil
}

func SearchUsersByPreference(data map[string]any, users *[]User) error {
	err := HelixClient.Query(
		"search_users_by_preference",
		helix.WithData(data),
	).Scan(helix.WithDest("users", users))
	if err != nil {
		err := fmt.Errorf("Error while searching for users by preferences: %v", err)
		return err
	}

	return nil
}
