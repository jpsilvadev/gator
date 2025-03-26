package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User set to:", cmd.Args[0])
	return nil
}
