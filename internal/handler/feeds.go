package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

func HandleAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]
	ctx := context.Background()

	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	follow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Printf("Feed created: %+v\n", feed)
	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}

func HandleFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s | URL: %s | User: %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}
