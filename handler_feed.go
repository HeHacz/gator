package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hehacz/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf(("agg cmd expects a one argument, time between reqests"))
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time format or duratiation error: %v", err)
	}
	fmt.Printf("Colletction feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf(("addfeed cmd expects a two arguments, the name and the feed URL"))
	}
	name := cmd.args[0]
	feedURL := cmd.args[1]
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
	ffpayload := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	}
	ffollow, err := s.db.CreateFeedFollow(context.Background(), ffpayload)
	if err != nil {
		return fmt.Errorf("coundn't create feed follow error: %v", err)
	}
	fmt.Printf("Feed \"%s\" with URL \"%s\" has been successfuly created: %s\n", feed.Name, feed.Url, user.Name)
	fmt.Printf("Feed has been successfuly followed:\n")
	printFeedFollow(ffollow.UserName, ffollow.FeedName)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
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

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch feed from database error: %v", err)
	}
	fmt.Printf("Found feed to fetch!\n")
	if _, err := s.db.MarkFeedFetched(context.Background(), next_feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed %s as fetched error: %v", next_feed.Name, err)
	}
	scrapeFeed(next_feed)
	return nil
}

func scrapeFeed(feed database.Feed) error {
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed data error: %v", err)
	}
	for i, item := range feedData.Channel.Item {
		fmt.Printf("%d. post found: %s\n", i, item.Title)
	}
	fmt.Printf("Feed %s collected. Found %d posts\n", feed.Name, len(feedData.Channel.Item))
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
