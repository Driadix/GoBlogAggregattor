package cli

import (
	"context"
	"fmt"
	"gator/internal/database"
	"strconv"
)

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := int32(2)
	if len(cmd.Args) > 0 {
		if parsed, err := strconv.ParseInt(cmd.Args[0], 10, 32); err == nil {
			limit = int32(parsed)
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("an error occured while fetching posts: %w", err)
	}

	for _, post := range posts {
		fmt.Println(post)
	}
	return nil
}
