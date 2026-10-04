package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.registeredCommands == nil {
		c.registeredCommands = make(map[string]func(*state, command) error)
	}
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("%s command is not registered", cmd.Name)
	}
	return f(s, cmd)
}
