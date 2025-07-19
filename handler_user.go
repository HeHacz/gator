package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hehacz/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login cmd expects a single argument, the username")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("coudn't recive user %s from the database: %v", cmd.args[0], err)
	}
	if user.ID != uuid.Nil {
		if err := s.conf.SetUser(cmd.args[0]); err != nil {
			return fmt.Errorf("coudn't set user to %s: %v", cmd.args[0], err)
		}
		fmt.Printf("User has been set to %s\n", cmd.args[0])
	} else {
		fmt.Printf("User not found in the database\n")
		return fmt.Errorf("user not found in the database")
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register cmd expects a single argument, the name")
	}
	payload := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	user, err := s.db.CreateUser(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("coudn't register user %s: %v", cmd.args[0], err)
	}
	fmt.Printf("User %v has been registered\n", user)
	handlerLogin(s, cmd)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DropUsers(context.Background()); err != nil {
		return fmt.Errorf("error during db resset: %v", err)
	}
	fmt.Printf("Database has been successufly resetarted!!")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve any user from the database error: %v", err)
	}
	for _, user := range users {
		if user.Name == s.conf.Current_user_name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s", user.Name)
		}
	}
	return nil
}
