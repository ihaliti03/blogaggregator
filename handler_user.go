package main

import (
	"context"
	"fmt"
	"time"
	
	//"uuid"
	"github.com/google/uuid"
	"github.com/ihaliti03/blogaggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error{
	// don't check with cmd.arguments == nil but use len()
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argumments inputted")
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(cmd.arguments[0]);err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("usage: register <username>")
	}
	
	username := cmd.arguments[0]

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("user %s created\n", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err!=nil{
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err!=nil{return err}

	for _, user := range users{
		if s.cfg.CurrentUserName == user.Name{
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error{
	const feedURL = "https://www.wagslane.dev/index.xml"
	
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n",feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error{
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("usage: addfeed <name> <link>")
	}

	name := cmd.arguments[0]
	feedURL := cmd.arguments[1]

	ctx := context.Background()
	currentUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
    if err != nil {
        return fmt.Errorf("couldn't find current user %q: %w", s.cfg.CurrentUserName, err)
    }

	feedRequest := database.CreateFeedParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name:      name,
        Url:       feedURL,
        UserID:    currentUser.ID,
    }

	feed, err := s.db.CreateFeed(ctx, feedRequest)
    if err != nil {
        return fmt.Errorf("couldn't create feed: %w", err)
    }

	fmt.Printf("Feed created: %+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error{
	feeds, err := s.db.GetFeeds(context.Background())
	if err!=nil{
		return err
	}

	for _, feed := range feeds{
		fmt.Printf("Name: %s | URL: %s | User: %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}