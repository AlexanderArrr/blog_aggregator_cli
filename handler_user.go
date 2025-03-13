package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("handlerLogin() expected 1 argument, received %v", len(cmd.args))
	}

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		os.Exit(1)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("The Username %v has been set!\n", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("handlerRegister expected 1 argument, received %v", len(cmd.args))
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	_, err := s.db.GetUser(context.Background(), userParams.Name)
	if err == nil {
		os.Exit(1)
	}
	_, err = s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("A user with the name %v was created!\n", cmd.args[0])
	fmt.Printf("User ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\n",
		userParams.ID, userParams.CreatedAt, userParams.UpdatedAt, userParams.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DelUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	userList, err := s.db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, user := range userList {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}
