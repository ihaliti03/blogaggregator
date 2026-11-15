package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.arguments) < 2 {
        return fmt.Errorf("usage: addfeed <name> <url>")
    }

    name := cmd.arguments[0]
    feedURL := cmd.arguments[1]
    ctx := context.Background()
	
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
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

    follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
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

func handlerFeeds(s *state, cmd command) error{
	feeds, err := s.db.GetFeeds(context.Background())
	if err!=nil{
		return err
	}

	for _, feed := range feeds{
		fmt.Printf("Name: %s | URL: %s | User: %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}