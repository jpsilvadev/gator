package main

import (
	"context"
	"fmt"

	"github.com/jpsilvadev/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %v", err)
		}

		if err := handler(s, cmd, user); err != nil {
			return fmt.Errorf("failed to handle command: %v", err)
		}
		return nil
	}
}
