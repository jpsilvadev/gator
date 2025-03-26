package main

import (
	"fmt"
	"os"

	"github.com/jpsilvadev/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	gatorState := &state{
		config: &cfg,
	}

	cmds := commands{
		cmdToHandler: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	argtoCmd := os.Args[2:]

	err = cmds.run(gatorState, command{
		Name: cmd,
		Args: argtoCmd,
	})
	if err != nil {
		fmt.Printf("Error executing %s: %v\n", cmd, err)
		os.Exit(1)
	}
}
