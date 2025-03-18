package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: gator addfeed <name> <url>")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("error while creating feed: %v", err)
	}

	var arguments []string
	arguments = append(arguments, feedURL)
	followdCmd := command{
		name: "follow",
		args: arguments,
	}
	err = handlerFollow(s, followdCmd, user)
	if err != nil {
		return fmt.Errorf("error while following new feed: %v", err)
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: gator feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error while getting feeds: %v", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error while getting user by id: %v", err)
		}
		fmt.Printf("Feed name: %v\n", feed.Name)
		fmt.Printf("Feed URL: %v\n", feed.Url)
		fmt.Printf("Added by: %v\n", user.Name)
	}

	return nil
}
