package internal

// Contains `Follow` related queries

import (
	"fmt"

	"github.com/HelixDB/helix-go"
)

type FollowUserInput struct {
	FollowerId string `json:"followerId"`
	FollowedId string `json:"followedId"`
}

func FollowUser(data *FollowUserInput) error {
	_, err := HelixClient.Query(
		"follow",
		helix.WithData(data),
	).Raw()
	if err != nil {
		err = fmt.Errorf("Error while following: %s", err)
		return err
	}

	return nil
}

func Followers(data map[string]any, users *[]User) error {
	err := HelixClient.Query(
		"followers",
		helix.WithData(data),
	).Scan(
		helix.WithDest("followers", users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"followers\": %s", err)
		return err
	}

	return nil
}

func Following(data map[string]any, users *[]User) error {
	err := HelixClient.Query(
		"following",
		helix.WithData(data),
	).Scan(
		helix.WithDest("following", users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"following\": %s", err)
		return err
	}

	return nil
}
