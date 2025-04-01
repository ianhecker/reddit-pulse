package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/ianhecker/reddit-pulse/config"
	"github.com/ianhecker/reddit-pulse/errorChecker"
	"github.com/ianhecker/reddit-pulse/logger"
	"github.com/ianhecker/reddit-pulse/poller"
)

const Verbose bool = true
const SleepDuration = 1 * time.Second
const PostCount = 100
const OutputFileName = "output.json"

type Output struct {
	Posts      poller.Posts
	TopAuthors []*poller.Author
}

func main() {
	ec := errorChecker.NewErrorChecker()
	log := logger.MakeLogger()
	log.SetVerbose(Verbose)

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

	postPoller, err := poller.NewPoller(credentials, userAgent, tokenURL)
	ec.WithMessage("could not create poller").CheckErr(err)

	ctx := context.Background()
	for {
		poll := postPoller.TopPosts(ctx, "golang", PostCount)
		ec.WithMessage("error polling").CheckErr(poll.Error)
		log.Log("polled top %d posts", PostCount)

		mostPosts := poller.MakeAuthors()
		mostPosts.CountPosts(poll.Posts)
		log.Log("counted posts for each author")

		topAuthors := mostPosts.TopAuthorsForCount(PostCount)
		log.Log("took top %d authors", PostCount)

		out := Output{Posts: poll.Posts, TopAuthors: topAuthors}

		bytes, err := json.MarshalIndent(out, "", " ")
		ec.WithMessage("could not marshal posts").CheckErr(err)
		log.Log("marshaled posts & authors to JSON")

		err = os.WriteFile(OutputFileName, bytes, 0644)
		ec.WithMessage("could not write to file").CheckErr(err)
		log.Log("wrote to file: %s", OutputFileName)

		remaining, used, seconds, err := poll.Response.GetRateLimits()
		ec.WithMessage("error converting X-RateLimit from header").CheckErr(err)
		log.Log("rate limits: requests remaining: %d, used: %d, seconds left: %d", remaining, used, seconds)

		log.Log("sleeping for duration: %s", SleepDuration)
		time.Sleep(SleepDuration)
	}
}
