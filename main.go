package main

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/ianhecker/reddit-pulse/config"
	"github.com/ianhecker/reddit-pulse/errorChecker"
	"github.com/ianhecker/reddit-pulse/poller"
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

	poller, err := poller.NewPoller(credentials, userAgent, tokenURL)
	ec.WithMessage("could not create poller").CheckErr(err)

	ctx := context.Background()

	posts, response, err := poller.TopPosts(ctx, "golang", 3)
	ec.WithMessage("error fetching posts").CheckErr(err)

	remaining, err := response.RequestsRemaining()
	ec.WithMessage("error getting remaining requests").CheckErr(err)

	requests, err := response.RequestsUsed()
	ec.WithMessage("error getting total requests").CheckErr(err)

	seconds, err := response.SecondsUntilReset()
	ec.WithMessage("error getting remaining seconds").CheckErr(err)

	fmt.Printf("Remaining Requests: %d\n", remaining)
	fmt.Printf("Total Requests Used: %d\n", requests)
	fmt.Printf("Seconds unti Reset: %d\n", seconds)

	for _, post := range posts {
		fmt.Printf("subreddit[%s] title:%s score:%d \nurl:%s\n", post.SubredditName, post.Title, post.Score, post.Permalink)
	}
}
