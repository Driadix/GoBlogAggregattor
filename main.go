package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registered[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registered[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register command expects a username argument")
	}

	newName := cmd.args[0]

	// 1. Create User in DB
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      newName,
	})
	if err != nil {
		return fmt.Errorf("error while creating a new user: %w", err)
	}

	// 2. FIX: Persist to config file!
	// Just setting s.cfg.CurrentUserName = newName is strictly in-memory.
	if err := s.cfg.SetUser(newName); err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	fmt.Printf("New user was created: %v\n", newUser)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login command expects a username argument")
	}

	user := cmd.args[0]

	// 1. Verify user exists in DB
	_, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("no user found: %w", err)
	}

	// 2. FIX: Actually set the user in the config!
	// You removed this line, but it is required for 'login' to work.
	if err := s.cfg.SetUser(user); err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	fmt.Printf("User set to: %s\n", user)
	return nil
}

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

	appState := &state{
		cfg: cfg,
		db:  dbQueries,
	}

	cmds := &commands{
		registered: make(map[string]func(*state, command) error),
	}

	// 3. FIX: Register the new command!
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err := cmds.run(appState, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
