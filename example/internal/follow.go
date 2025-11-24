package internal

// Contains `Follow` related queries

import (
	"fmt"

	"github.com/HelixDB/helix-go"
)

type FollowUserInput struct {
	FollowerId string `json:"follower_id"`
	FollowedId string `json:"followed_id"`
}

func FollowUser(data *FollowUserInput) error {
	_, err := HelixClient.Query(
		"follow",
		helix.WithData(data),
	)
	if err != nil {
		err = fmt.Errorf("Error while following: %v", err)
		return err
	}

	return nil
}

func Followers(data map[string]any, users *[]User) error {
	res, err := HelixClient.Query(
		"followers",
		helix.WithData(data),
	)
	if err != nil {
		return fmt.Errorf("Error while getting \"followers\": %v", err)
	}

	if err = res.Scan(helix.WithDest("followers", users)); err != nil {
		return fmt.Errorf("Error while scanning the \"followers\" query result and populating the users slice: %v", err)
	}

	return nil
}

func Following(data map[string]any, users *[]User) error {
	res, err := HelixClient.Query(
		"following",
		helix.WithData(data),
	)
	if err != nil {
		return fmt.Errorf("Error while getting \"following\": %v", err)
	}

	if err = res.Scan(helix.WithDest("following", users)); err != nil {
		return fmt.Errorf("Error while scanning the \"following\" query result and populating the users slice: %v", err)
	}

	return nil
}
