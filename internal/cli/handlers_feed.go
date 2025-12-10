package cli

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("the addfeed command expects two arguments: <name> <url>")
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	newFeedID := uuid.New()

	newFeed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        newFeedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("cannot create a new feed: %w", err)
	}

	_, err = s.DB.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    newFeedID,
		},
	)

	fmt.Println(newFeed)
	return nil
}

func HandlerFeeds(s *State, _ Command) error {
	feeds, err := s.DB.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  User: %s\n", feed.UserName)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the folllow command expects argument: <url>")
	}

	foundFeed, err := s.DB.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("can't find feed with such url, err: %w", err)
	}

	followedFeed, err := s.DB.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    foundFeed.ID,
		},
	)

	fmt.Println(followedFeed)
	return nil
}

func HandlerFollowing(s *State, _ Command, user database.User) error {
	followingFeeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error while fetching follows for user: %w", err)
	}

	if len(followingFeeds) == 0 {
		fmt.Println("No following feeds found")
		return nil
	}

	for _, feed := range followingFeeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
