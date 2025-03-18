package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator agg <duration-string> (for example: '1s' or '1m' or '1h')")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error while parsing duration-string: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error while getting next feed to fetch %v", err)
	}

	feedToFetch, err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return fmt.Errorf("error while marking feed as fetched: %v", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("error while getting feed by url: %v", err)
	}

	for _, feed := range rssFeed.Channel.Item {
		const longForm = "Fri, 07 Mar 2025 00:00:00 +0000"
		publishedAt, err := time.Parse(longForm, feed.PubDate)
		if err != nil {
			return fmt.Errorf("error while parsing publish date: %v", err)
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       feed.Title,
			Url:         feed.Link,
			Description: feed.Description,
			PublishedAt: publishedAt,
			FeedID:      feedToFetch.ID,
		}
		_, err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			return fmt.Errorf("error while creating post: %v", err)
		}
	}

	fmt.Printf("Saved %v posts from the feed %v!", len(rssFeed.Channel.Item), rssFeed.Channel.Title)

	return nil
}
