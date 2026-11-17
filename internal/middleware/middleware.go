package middleware

import (
	"context"
	"fmt"

	"github.com/ihaliti03/blogaggregator/internal/database"
	"github.com/ihaliti03/blogaggregator/internal/handler"
)

func LoggedIn(next func(s *handler.State, cmd handler.Command, user database.User) error) func(*handler.State, handler.Command) error {
	// ═══ MOMENT 1: runs now (at registration) ═══
	// `handler` is captured.
	return func(s *handler.State, cmd handler.Command) error {
		// ═══ MOMENT 2: runs later (at dispatch) ═══
		// `handler` is available here (captured).
		// `GetUser` runs here.
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find current user %q — run `blogaggregator login` first: %w", s.Cfg.CurrentUserName, err)
		}
		return next(s, cmd, user)
	}
}
