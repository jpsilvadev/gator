package main

import (
	"context"
	"fmt"
	"log"
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

	logFeedInfo(feed)
	return nil
}

func logFeedInfo(feed database.Feed) {
	log.Printf("Feed created:\n  ID: %s\n  CreatedAt: %s\n  UpdatedAt: %s\n  Name: %s\n  URL: %s\n  UserID: %s\n",
		feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
}
