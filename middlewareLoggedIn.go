package main

import (
	"context"
	"fmt"

	"github.com/hehacz/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}
		return handler(s, cmd, user)
	}
}
