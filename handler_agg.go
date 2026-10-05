package main

import (
	"context"
	"fmt"
	"log"
	"time"
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

	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Description:\n%s\n", item.Description)
	}
	log.Printf("Feed %s collected, %v posts found", nextFeedFetch.Name, len(feed.Channel.Item))
	return nil
}
