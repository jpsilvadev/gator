package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	fmt.Printf("User has been set to %s\n", name)
	return nil
}
