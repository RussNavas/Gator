package main
import (
	"Gator/internal/database"
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error{
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: follow <url>")
	}
	url := cmd.Args[1]
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
		return fmt.Errorf("error creating FeedFollow record")
	}
	fmt.Println("feed_follow record created successfully!")
	fmt.Printf("Feed Name: %v\n", feedFollowRow.FeedName)
	fmt.Printf("Current User: %v\n", currUser.Name)
	fmt.Println("====================================")
	return nil
}
