package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: gator browse <limit> (limit is optional!)")
	}

	limit := 2
	if len(cmd.args) == 1 {
		limit, _ = strconv.Atoi(cmd.args[0])
	}

	getPostsParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostsParams)
	if err != nil {
		return fmt.Errorf("error while getting posts for user: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Post title: %v", post.Title)
	}

	return nil
}
