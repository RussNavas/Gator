package main
import (
	"Gator/internal/database"
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error{

	url := cmd.Args[0]
	currUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user")
	}
	currUserId := currUser.ID

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil{
		return fmt.Errorf("error getting feed from url")
	}
	feedID := feed.ID

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: currUserId,
		FeedID: feedID,
	})
	if err != nil {
		return fmt.Errorf("error creating FeedFollow record during add feed: %v", err)
	}
	fmt.Println("feed_follow record created successfully!")
	fmt.Printf("Feed Name: %v\n", feedFollowRow.FeedName)
	fmt.Printf("Current User: %v\n", currUser.Name)
	fmt.Println("====================================")
	return nil
}

func handlerFollowing(s *state, cmd command) error {

	user := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), user)
	if err != nil{
		return fmt.Errorf("error getting current user")
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil{
		return fmt.Errorf("eror getting feed follows for user %+v", currentUser)
	}
	fmt.Println("Printing followed feeds:")
	for _, feed := range feeds{
		fmt.Printf("feed name: %v\n", feed.FeedName)
	}
	fmt.Println("======================================")
	return nil
}
