package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, fmt.Errorf("GET %s: status %d: %s", feedURL, res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("unmarshalling xml: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	if err := s.db.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed fetched: %w", err)
	}

	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %q: %w", feed.Url, err)
	}

	for _, item := range rss.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		publishedAt := sql.NullTime{}
		if item.PubDate != "" {
			if t, err := parsePubDate(item.PubDate); err == nil {
				publishedAt = sql.NullTime{Time: t, Valid: true}
			} else {
				fmt.Printf("couldn't parse pubDate %q: %v\n", item.PubDate, err)
			}
		}

		err := s.db.CreatePost(ctx, database.CreatePostParams{
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

	fmt.Printf("Feed: %s — saved %d items\n", feed.Name, len(rss.Channel.Item))

	return nil
}

func parsePubDate(s string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"  ← most common
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		time.RFC3339,  // "2006-01-02T15:04:05Z07:00"
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("no matching layout for %q", s)
}
