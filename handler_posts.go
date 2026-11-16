package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ihaliti03/blogaggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {

	limit := int32(2)
	if len(cmd.arguments) > 0 {
		parsedInt, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return fmt.Errorf("error parsing limit, it needs to be integer: %w", err)
		}
		limit = int32(parsedInt)
	}

	ctx := context.Background()

	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user - %s: %w", user.Name, err)
	}

	if len(posts) == 0 {
		fmt.Println("no posts found — try running `./blogaggregator agg` first")
		return nil
	}

	fmt.Printf("Posts for %s:\n", user.Name)
	for _, p := range posts {
		fmt.Printf("\n* %s\n", p.Title)
		fmt.Printf("  Feed: %s\n", p.FeedName)
		fmt.Printf("  URL:  %s\n", p.Url)
		if p.PublishedAt.Valid {
			fmt.Printf("  Published: %s\n", p.PublishedAt.Time.Format("Mon Jan 2, 2006"))
		}
		if p.Description.Valid && p.Description.String != "" {
			desc := p.Description.String
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			fmt.Printf("  %s\n", desc)
		}
	}

	return nil
}
