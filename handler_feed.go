package main

import (
	"Gator/internal/database"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
	"github.com/google/uuid"
)


type RSSFeed struct {
	Channel struct{
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description	string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string 	`xml:"link"`
	Description	string 	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){

	rssFeed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context failed: %v", err)
	}

	req.Header.Set("User-Agent", "gator")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil{
		return nil, fmt.Errorf("response failed: err -> %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with the following status code: %v", res.StatusCode)
	}


	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read res body: %v", err)
	}

	defer res.Body.Close()

	err = xml.Unmarshal(respBody, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("unmarshal resp failed: %v", err)
	}

	// unescape channel (title & desc)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	// unescape item(s) (title & desc)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}


func handlerAgg(s *state, cmd command) error{
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <parse duration w/ unit of time>")
	}

	timeToParse, err := time.ParseDuration(cmd.Args[0])
	if err != nil{
		return fmt.Errorf("unable to parse time for valid duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeToParse)
	ticker := time.NewTicker(timeToParse)
	for ; ; <- ticker.C {
		err = scrapeFeeds((s))
		if err != nil {
			fmt.Printf("error scraping feeds: %v", err)
		}
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error{
	if len(cmd.Args) != 2 {
		return fmt.Errorf("not enough args, usage <name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	fmt.Println("Getting current user ...")

	fmt.Println("Creating feed ...")
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:			uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})

	if err != nil{
		return fmt.Errorf("problem creating feed in db: %v", err)
	}

	fmt.Println("Creating FeedFollow ...")

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating feedfollow for user %v", err)
	}

	fmt.Println("Feed created sucessfully:")
	printFeed(feed)
	fmt.Println("=========================================")

	return nil
}


func printFeed(feed database.Feed){
	fmt.Printf("ID:  		%v\n", feed.ID)
	fmt.Printf("CreatedAt: 	%v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: 	%v\n", feed.UpdatedAt)
	fmt.Printf("Name: 		%v\n", feed.Name)
	fmt.Printf("Url: 		%v\n", feed.Url)
	fmt.Printf("UserID: 	%v\n", feed.UserID)
}

func handlerFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil{
		return fmt.Errorf("error getting feeds from db: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		printFeed(feed)
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user by id: %v", err)
		}
		fmt.Printf("UserName: 	%v\n", user.Name)
		fmt.Println("================================")
	}

	return nil
}

func scrapeFeeds(s *state) error {
	fmt.Println("Scraping feeds ...")
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return fmt.Errorf("unable to GetNextFeedToFetch")
	}

	feedRSS, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to get feedRSS by URL")
	}

	feed, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("unable to MarkFeedFetched")
	}

	fmt.Println("=============== Feeds ==============")
	for i, item := range feedRSS.Channel.Item {
		fmt.Printf("|------Printing Feed Title %v ----------|\n", i)
		fmt.Print("")
		fmt.Printf("%v\n", item.Title)
		fmt.Print("")
		fmt.Println("----------------------------------")
		fmt.Print("")
	}
	return nil
}
