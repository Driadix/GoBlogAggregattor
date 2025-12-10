package middleware

import (
	"context"
	"fmt"
	cli "gator/internal/cli"
	"gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *cli.State, cmd cli.Command, user database.User) error) func(*cli.State, cli.Command) error {
	return func(s *cli.State, cmd cli.Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("authentication error: %w", err)
		}
		return handler(s, cmd, user)
	}
}
