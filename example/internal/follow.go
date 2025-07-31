package internal

// Contains `Follow` related queries

import (
	"fmt"

	"github.com/HelixDB/helix-go"
)

func FollowUser(followerId string, followedId int) error {
	_, err := HelixClient.Query(
		"follow",
		helix.WithData(followerId),
		helix.WithData(followedId),
	).Raw()
	if err != nil {
		err = fmt.Errorf("Error while getting users: %s", err)
		return err
	}

	return nil
}

func Followers(userId string, users *[]User) error {
	err := HelixClient.Query(
		"followers",
		helix.WithData(userId),
	).Scan(
		helix.WithDest("followers", users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"followers\": %s", err)
		return err
	}

	return nil
}

func FollowerCount(userId string) (*int, error) {
	var count int

	err := HelixClient.Query(
		"followerCount",
		helix.WithData(userId),
	).Scan(
		helix.WithDest("count", &count),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"follower count\": %s", err)
		return nil, err
	}

	return &count, nil
}

func Following(userId string, users *[]User) error {
	err := HelixClient.Query(
		"following",
		helix.WithData(userId),
	).Scan(
		helix.WithDest("following", users),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"following\": %s", err)
		return err
	}

	return nil
}

func FollowingCount(userId string) (*int, error) {
	var count int

	err := HelixClient.Query(
		"followingCount",
		helix.WithData(userId),
	).Scan(
		helix.WithDest("count", &count),
	)
	if err != nil {
		err = fmt.Errorf("Error while getting \"following count\": %s", err)
		return nil, err
	}

	return &count, nil
}
