package cli

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while fetching feed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("error while parsing xml: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	scrapeFeeds(s)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *State) error {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	err = s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed %s fetched: %w", feed.Name, err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("Couldn't collect feed %s: %v\n", feed.Name, err)
		return nil
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))

	return nil
}
