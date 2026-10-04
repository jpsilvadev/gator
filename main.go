package main

import (
	"log"
	"os"

	"github.com/jpsilvadev/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	gatorState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := args[1]
	cmdArgs := args[2:]
	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}
	err = cmds.run(gatorState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
