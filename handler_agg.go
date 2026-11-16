package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	// beautiful use of time.Parse
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.arguments[0], err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feeds: %v\n", err)
		}
	}
}
