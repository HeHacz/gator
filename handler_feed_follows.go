package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hehacz/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf(("follow cmd expects one argument, the feed URL"))
	}
	feed, err := s.db.GetFeedbyURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("coudnt get feed with url %s error: %v", cmd.args[0], err)
	}
	userID, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		return fmt.Errorf("couldn't retrieve current user id from database error: %v", err)
	}
	payload := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID.ID,
		FeedID:    feed.ID,
	}
	ffollow, err := s.db.CreateFeedFollow(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("coundn't create feed follow error: %v", err)
	}
	fmt.Printf("Feed follow created: \n")
	printFeedFollow(ffollow.UserName, ffollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	userID, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		return fmt.Errorf("couldn't retrieve current user id from database error: %v", err)
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), userID.ID)
	if err != nil {
		return fmt.Errorf("couldn't retrieve current user feeds from the database error: %v", err)
	}
	for _, feed := range feeds {
		printFeedFollow(s.conf.Current_user_name, feed.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:		%s\n", username)
	fmt.Printf("* Feed: 	%s\n", feedname)
}
