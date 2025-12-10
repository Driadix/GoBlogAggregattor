package cli

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	registered map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		registered: make(map[string]func(*State, Command) error),
	}
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.registered[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.registered[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return f(s, cmd)
}
