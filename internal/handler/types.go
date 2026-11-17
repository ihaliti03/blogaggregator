package handler

import (
	"github.com/ihaliti03/blogaggregator/internal/config"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}

type Command struct {
	Name      string
	Arguments []string
}

type HandlerFunc func(s *State, cmd Command) error
type LoggedInFunc func(s *State, cmd Command, user database.User) error
