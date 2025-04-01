package main

import (
	"context"
	"fmt"
	"time"

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

	for {
		poll := poller.TopPosts(ctx, "golang", 1)
		if poll.Error != nil {
			fmt.Errorf("error fetching top posts: %s", err)
			continue
		}

		bytes, err := poll.Posts.MarshalJSON()
		ec.WithMessage("could not marshal posts").CheckErr(err)

		fmt.Println(string(bytes))

		time.Sleep(1 * time.Second)
	}

	// ec.WithMessage("error fetching posts").CheckErr(err)

	// remaining, used, seconds, err := response.GetRateLimits()
	// ec.WithMessage("error getting rate limits").CheckErr(err)

	// fmt.Printf("Remaining Requests: %d\n", remaining)
	// fmt.Printf("Total Requests Used: %d\n", used)
	// fmt.Printf("Seconds unti Reset: %d\n", seconds)

	// for _, post := range posts {
	// 	fmt.Printf("subreddit[%s] title:%s score:%d \nurl:%s\n", post.SubredditName, post.Title, post.Score, post.Permalink)
	// }
}
