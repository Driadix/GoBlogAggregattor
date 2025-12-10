package main

import (
	"database/sql"
	"fmt"
	"gator/internal/cli"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to db: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	appState := &cli.State{
		Cfg: cfg,
		DB:  dbQueries,
	}

	cmds := cli.NewCommands()
	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)
	cmds.Register("users", cli.HandlerUsers)
	cmds.Register("agg", cli.HandlerAgg)
	cmds.Register("addfeed", cli.HandlerAddFeed)
	cmds.Register("feeds", cli.HandlerFeeds)
	cmds.Register("follow", cli.HandlerFollow)
	cmds.Register("following", cli.HandlerFollowing)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(appState, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
