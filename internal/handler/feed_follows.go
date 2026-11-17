package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

func HandleFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	url := cmd.Arguments[0]
	ctx := context.Background()

	feed, err := s.Db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't find feed with url %q: %w", url, err)
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

	fmt.Printf("%s is now following %s\n", follow.UserName, follow.FeedName)
	return nil
}

func HandleFollowing(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	follows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("no feeds followed")
		return nil
	}

	fmt.Printf("%s is following:\n", user.Name)
	for _, f := range follows {
		fmt.Printf(" - %s\n", f.FeedName)
	}
	return nil
}

func HandleUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	url := cmd.Arguments[0]
	ctx := context.Background()

	feed, err := s.Db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't find feed with url %q: %w", url, err)
	}

	if err := s.Db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("Feed: %s successfully unfollowed by user: %s\n", url, user.Name)
	return nil
}
