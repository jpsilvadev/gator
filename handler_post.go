package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jpsilvadev/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	// arbitrary default
	limit := 2
	if len(cmd.Args) == 1 {
		// assume we got the optional limit arg and no more
		specifiedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("failed parsing optional [limit] arg %q: %w", cmd.Args[0], err)
		}
		limit = specifiedLimit
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}

	getPostParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return fmt.Errorf("could not get posts for user %v: %w", user.Name, err)
	}

	for _, post := range posts {
		logPostInfo(post)
	}

	return nil
}

func logPostInfo(post database.GetPostsForUserRow) {
	fmt.Println("==================================================")
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Source: %v\n", post.FeedName)
	fmt.Printf("Published at: %v\n", post.PublishedAt.Time.Format(time.DateTime))

	fmt.Printf("Link: %s\n", post.Url)
	fmt.Printf("Description: %v\n", post.Description.String)
}
