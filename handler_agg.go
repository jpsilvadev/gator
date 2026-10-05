package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jpsilvadev/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse <time_between_reqs> duration string: %w\n", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Printf("error scraping feed: %v", err)
			continue
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeedFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed to fetch: %w\n", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeedFetch.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed %s fetched: %w", nextFeedFetch.ID, err)
	}

	feed, err := fetchFeed(context.Background(), nextFeedFetch.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	newPosts := 0
	for _, item := range feed.Channel.Item {
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: parseDescription(item.Description),
			PublishedAt: parsePubDate(item.PubDate),
			FeedID:      nextFeedFetch.ID,
		}
		_, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				// duplicated posts are expected
				// should not crash the scraper
				continue
			}
			log.Printf("could not create post: %v", err)
			continue
		}
		newPosts++
	}
	log.Printf("Feed %s collected, %v posts found, %v new", nextFeedFetch.Name, len(feed.Channel.Item), newPosts)
	return nil
}

// cast post description to nullable in case its empty
func parseDescription(description string) sql.NullString {
	parsedDescription := sql.NullString{
		String: description,
		Valid:  description != "",
	}
	return parsedDescription
}

// cast post publication date to nullable in case its empty
func parsePubDate(pubDate string) sql.NullTime {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		// single-digit day variants, e.g. "Mon, 2 Jan 2006"
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, pubDate); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}
	return sql.NullTime{Valid: false}
}
