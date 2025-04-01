package main

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/ianhecker/reddit-pulse/config"
	"github.com/ianhecker/reddit-pulse/errorChecker"
)

func main() {
	ec := errorChecker.NewErrorChecker()

	cfg, err := config.NewConfig()
	ec.WithMessage("could not make config").CheckErr(err)

	credentials := reddit.Credentials{
		ID:       cfg.ClientID,
		Secret:   cfg.ClientSecret,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	userAgent := reddit.WithUserAgent(cfg.UserAgent)
	tokenURL := reddit.WithTokenURL("https://www.reddit.com/api/v1/access_token")

	client, err := reddit.NewClient(credentials, userAgent, tokenURL)
	ec.WithMessage("could not create reddit client").CheckErr(err)

	ctx := context.Background()

	posts, _, err := client.Subreddit.TopPosts(ctx, "golang", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 100,
		},
		Time: "all",
	})
	ec.WithMessage("error fetching posts").CheckErr(err)

	for _, post := range posts {
		fmt.Printf("[%s] %s\n", post.SubredditName, post.Title)
	}
}
