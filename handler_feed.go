package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jpsilvadev/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := s.cfg.CurrentUserName

	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}

	name = cmd.Args[0]
	url := cmd.Args[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}
	fmt.Println("Feed created:")
	logFeedInfo(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserInfoByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("could not get user info: %w", err)
		}
		fmt.Printf("==== User: %s ====\n", user.Name)
		logFeedInfo(feed)
		fmt.Println("==================================================")

	}
	return nil
}

func logFeedInfo(feed database.Feed) {
	fmt.Printf("ID: %s\n  CreatedAt: %s\n  UpdatedAt: %s\n  Name: %s\n  URL: %s\n  UserID: %s\n",
		feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
}
