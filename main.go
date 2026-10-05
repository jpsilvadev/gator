package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jpsilvadev/gator/internal/config"
	"github.com/jpsilvadev/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// load config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// load db
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	// setup global state
	gatorState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// setup commands
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("reset", handlerReset)
	cmds.register("feeds", handlerListFeeds)

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
