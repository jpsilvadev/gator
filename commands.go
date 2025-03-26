package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmdToHandler map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdToHandler[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmdToHandler[cmd.Name]
	if !ok {
		return errors.New("unknown command: " + cmd.Name)
	}
	return f(s, cmd)
}
