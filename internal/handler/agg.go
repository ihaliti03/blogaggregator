package handler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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

const maxConcurrentFetches = 10

func ScrapeFeeds(s *State) error {
	ctx := context.Background()

	feeds, err := s.Db.GetFeedsToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get feeds to fetch: %w", err)
	}

	if len(feeds) == 0 {
		return nil
	}

	// semaphore: at most maxConcurrentFetches goroutines run at once
	sem := make(chan struct{}, maxConcurrentFetches)
	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		sem <- struct{}{} // acquire a slot (blocks if 10 already running)
		go func(f database.Feed) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot

			if err := fetchAndSave(s, &f); err != nil {
				fmt.Printf("error fetching %q: %v\n", f.Name, err)
			}
		}(feed)
	}

	wg.Wait()
	return nil
}

func fetchAndSave(s *State, feed *database.Feed) error {
	ctx := context.Background()

	// mark first (so it's not re-fetched next tick)
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
