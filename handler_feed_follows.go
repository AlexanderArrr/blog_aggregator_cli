package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator follow <url>")
	}
	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error while getting feed by URL: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error while creating FeedFollow record: %v", err)
	}
	fmt.Printf("Feed name: %v\n", feedFollow.FeedName)
	fmt.Printf("User name: %v\n", feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: gator following")
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error while getting feed follows for user: %v", err)
	}

	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}

	return nil
}
