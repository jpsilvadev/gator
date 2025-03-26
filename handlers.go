package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jpsilvadev/gator/internal/database"
	"github.com/jpsilvadev/gator/internal/rss"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// check if user exists
	user, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("unexpected database error: %v", err)
	}
	if err == nil {
		return fmt.Errorf("user %s already exists", user.Name)
	}

	// create user
	user, err = s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user %s", cmd.Args[0])
	}

	// set active user
	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user %w", err)
	}
	fmt.Println("User created:", cmd.Args[0])

	// log user info TODO: remove later
	fmt.Println("User info:", user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == sql.ErrNoRows {
		return fmt.Errorf("user %s not found", cmd.Args[0])
	}
	if err != nil {
		return fmt.Errorf("unexpected database error: %v", err)
	}

	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User set to:", cmd.Args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %v", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list users: %v", err)
	}
	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", u.Name)
			continue
		}
		fmt.Printf("* %v\n", u.Name)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}
	fmt.Println(feed)
	return nil
}
