package main
import (
	"Gator/internal/database"
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error{

	if len(cmd.Args) != 1{
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil{
		return fmt.Errorf("error getting feed from url")
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating FeedFollow record during add feed: %v", err)
	}
	fmt.Println("feed_follow record created successfully!")
	fmt.Printf("Feed Name: %v\n", feedFollowRow.FeedName)
	fmt.Printf("Current User: %v\n", user.Name)
	fmt.Println("====================================")
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return fmt.Errorf("eror getting feed follows for user %+v", user)
	}
	fmt.Println("Printing followed feeds:")
	for _, feed := range feeds{
		fmt.Printf("feed name: %v\n", feed.FeedName)
	}
	fmt.Println("======================================")
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1{
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	urlToDelete := cmd.Args[0]

	fmt.Printf("Preparing to unfollow: %s\n", urlToDelete)
	fmt.Println("Getting feed by URL ...")


	feed, err := s.db.GetFeedByURL(context.Background(), urlToDelete)
	if err != nil{
		return fmt.Errorf("error getting feed by url: %v", err)
	}

	fmt.Println("Initiating Unfollow Request ...")

	err = s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil{
		return fmt.Errorf("Error unfollowing given url: %v", err)
	}

	fmt.Println("Unfollow was successful!")
	fmt.Println("==============================")
	return nil
}
