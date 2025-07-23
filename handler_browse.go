package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hehacz/gator/internal/database"
)

func handledBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 1 {
		return fmt.Errorf("browse command expects at most one argument: the number of posts to print (default is 2)")
	}
	if len(cmd.args) == 1 {
		if userLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = userLimit
		} else {
			return fmt.Errorf("invalid limit format. error: %v", err)
		}
	}
	payload := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("coudnt retrive posts for user %s from database. error: %v", user.Name, err)
	}
	fmt.Printf("Posts collected. Found %d posts for user %s:\n", len(posts), user.Name)
	fmt.Println("==========================================")
	for i, post := range posts {
		fmt.Printf("From feed: %s\n", post.FeedName)
		fmt.Printf("%d. %s\n", i+1, post.Title)
		fmt.Printf("\tPublished at: %s\n", post.PublishedAt.Time.Format("Mon Jan 2 2006"))
		fmt.Printf("\tDescryption: %s\n", post.Description.String)
		fmt.Printf("\tURL: %s\n", post.Url)
		fmt.Println("------------------------------------------")
	}
	fmt.Println("==========================================")
	return nil
}
