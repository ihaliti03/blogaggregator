package handler

import (
	"context"
	"fmt"
	"time"

	//"uuid"
	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

func HandleLogin(s *State, cmd Command) error {
	// don't check with cmd.arguments == nil but use len()
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("no argumments inputted")
	}

	user, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return err
	}

	if err := s.Cfg.SetUser(cmd.Arguments[0]); err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", user.Name)
	return nil
}

func HandleRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: register <username>")
	}

	username := cmd.Arguments[0]

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("user %s created\n", user.Name)
	return nil
}

func HandleReset(s *State, cmd Command) error {
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func HandleUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if s.Cfg.CurrentUserName == user.Name {
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}
	return nil
}
