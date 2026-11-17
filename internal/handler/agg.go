package handler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
	"github.com/ihaliti03/blogaggregator/internal/rss"
)

func HandleAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	// beautiful use of time.Parse
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.Arguments[0], err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := ScrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feeds: %v\n", err)
		}
	}
}

func ScrapeFeeds(s *State) error {
	ctx := context.Background()

	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	if err := s.Db.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed fetched: %w", err)
	}

	parsedFeed, err := rss.Fetch(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %q: %w", feed.Url, err)
	}

	for _, item := range parsedFeed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		publishedAt := sql.NullTime{}
		if item.PubDate != "" {
			if t, err := rss.ParsePubDate(item.PubDate); err == nil {
				publishedAt = sql.NullTime{Time: t, Valid: true}
			} else {
				fmt.Printf("couldn't parse pubDate %q: %v\n", item.PubDate, err)
			}
		}

		err := s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("couldn't create post %q: %v\n", item.Title, err)
		}
	}

	fmt.Printf("Feed: %s — saved %d items\n", feed.Name, len(parsedFeed.Channel.Item))

	return nil
}
