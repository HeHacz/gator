package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hehacz/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf(("follow cmd expects one argument, the feed URL"))
	}
	feed, err := s.db.GetFeedbyURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("coudnt get feed with url %s error: %v", cmd.args[0], err)
	}

	payload := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't retrieve current user feeds from the database error: %v", err)
	}
	for _, feed := range feeds {
		printFeedFollow(s.conf.Current_user_name, feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf(("unfollow cmd expects one argument, the feed URL"))
	}
	feed, err := s.db.GetFeedbyURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("coudnt find feed with url %s error: %v", cmd.args[0], err)
	}
	payload := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFollow(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("coudnt delete follow for a user %s, with url: %s error: %v", user.Name, feed.Url, err)
	}
	fmt.Printf("Follow for a user %s, with url: %s has been successfuly removed", user.Name, feed.Url)
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:		%s\n", username)
	fmt.Printf("* Feed: 	%s\n", feedname)
}
