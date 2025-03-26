package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jpsilvadev/gator/internal/config"
	"github.com/jpsilvadev/gator/internal/database"
	"github.com/jpsilvadev/gator/internal/rss"
	_ "github.com/lib/pq" // Load the PostgreSQL driver
)

type state struct {
	config *config.Config
	db     *database.Queries
	rss    *rss.RSSFeed
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	gatorState := &state{
		config: &cfg,
		db:     dbQueries,
	}

	cmds := commands{
		cmdToHandler: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

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
