package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hehacz/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error during fetching feed: %v", err)
	}
	fmt.Printf("%v\n", rss)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf(("addfeed cmd expects a two arguments, the name and the feed URL"))
	}
	name := cmd.args[0]
	feedURL := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	userID := user.ID

	if name == "" || feedURL == "" {
		return fmt.Errorf("name and feed URL parameters cannot be empty")
	}
	payload := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedURL,
		UserID:    userID,
	}
	feed, err := s.db.CreateFeed(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("error during feed creation: %v", err)
	}
	fmt.Printf("Feed %s with URL %s has been added for user %s\n", feed.Name, feed.Url, feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUsername(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds from database")
	}
	if len(feeds) == 0 {
		fmt.Println("There is no feeds to print")
		return nil
	}
	fmt.Printf("\nFound %d feeds\n", len(feeds))
	fmt.Println("==========================================")
	for i, feed := range feeds {
		printFeeds(feed)
		if i+1 < len(feeds) {
			fmt.Println("------------------------------------------")
		}
	}
	fmt.Println("==========================================")
	return nil
}

func printFeeds(feed database.GetFeedsWithUsernameRow) {
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("Created At: %s\n", feed.CreatedAt)
	fmt.Printf("Updated At: %s\n", feed.UpdatedAt)
	fmt.Printf("Username: %s\n", feed.Username)
}
