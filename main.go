package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
		ec.WithMessage("error polling").CheckErr(poll.Error)

		fmt.Println("polled posts")

		bytes, err := json.MarshalIndent(poll.Posts, "", " ")
		ec.WithMessage("could not marshal posts").CheckErr(err)

		fmt.Println("marshaled posts")

		err = os.WriteFile("subreddit.json", bytes, 0644)
		ec.WithMessage("could not write to file").CheckErr(err)

		fmt.Println("wrote to file")

		time.Sleep(1 * time.Second)
	}
}
