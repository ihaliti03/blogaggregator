package main

import (
	"context"
	"fmt"

	"github.com/ihaliti03/blogaggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	// ═══ MOMENT 1: runs now (at registration) ═══
    // `handler` is captured. 
	return func(s *state, cmd command) error {
		// ═══ MOMENT 2: runs later (at dispatch) ═══
        // `handler` is available here (captured).
        // `GetUser` runs here.
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find current user %q — run `gator login` first: %w", s.cfg.CurrentUserName, err)
		}
		return handler(s, cmd, user)
	}
}