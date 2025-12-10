package cli

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login command expects a username argument")
	}

	user := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("no user found: %w", err)
	}

	if err := s.Cfg.SetUser(user); err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	fmt.Printf("User set to: %s\n", user)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register command expects a username argument")
	}

	newName := cmd.Args[0]

	newUser, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      newName,
	})
	if err != nil {
		return fmt.Errorf("error while creating a new user: %w", err)
	}

	if err := s.Cfg.SetUser(newName); err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	fmt.Printf("New user was created: %v\n", newUser)
	return nil
}

func HandlerReset(s *State, _ Command) error {
	err := s.DB.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("an error occured while trying to reset users: %w", err)
	}

	fmt.Printf("users were successfully reseted")
	return nil
}

func HandlerUsers(s *State, _ Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("can't load users from database: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		}
		fmt.Println(user.Name)
	}
	return nil
}
