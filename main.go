package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"

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

	subreddit := cfg.Subreddit

	postPoller, err := poller.NewPoller(credentials, userAgent, tokenURL)
	ec.WithMessage("could not create poller").CheckErr(err)

	ctx := context.Background()
	limiter := rate.NewLimiter(rate.Limit(1), 10)

	workerCount := 3
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {

			defer wg.Done()
			for {
				err := limiter.Wait(ctx)
				if err != nil {
					fmt.Printf("Worker %d: error waiting for limiter: %v\n", workerID, err)
					return
				}
				delay := Poll(ctx, workerID, postPoller, subreddit)
				log.Log("worker %d returned from polling", workerID)

				newLimit := rate.Every(delay)
				limiter.SetLimit(newLimit)
				log.Log("set new limiting rate to: %s", delay)
			}
		}(i)
	}
	wg.Wait()
}

func Poll(
	ctx context.Context,
	workerID int,
	postPoller *poller.Poller,
	subreddit string,
) time.Duration {
	ec := errorChecker.NewErrorChecker()
	log := logger.MakeLogger()
	log.SetVerbose(Verbose)

	log.Log("worker #%d started polling", workerID)

	poll := postPoller.TopPosts(ctx, subreddit, PostCount)
	ec.WithMessage("error polling").CheckErr(poll.Error)
	log.Log("----polled top %d posts for subreddit: %s", PostCount, subreddit)

	mostPosts := poller.MakeAuthors()
	mostPosts.CountPosts(poll.Posts)
	log.Log("----counted posts for each author")

	topAuthors := mostPosts.TopAuthorsForCount(PostCount)
	log.Log("----took top %d authors", PostCount)

	out := Output{Posts: poll.Posts, TopAuthors: topAuthors}

	bytes, err := json.MarshalIndent(out, "", " ")
	ec.WithMessage("could not marshal posts").CheckErr(err)
	log.Log("----marshaled posts & authors to JSON")

	err = os.WriteFile(OutputFileName, bytes, 0644)
	ec.WithMessage("could not write to file").CheckErr(err)
	log.Log("----wrote to file: %s", OutputFileName)

	remaining, used, seconds, err := poll.Response.GetRateLimits()
	ec.WithMessage("error converting X-RateLimit from header").CheckErr(err)
	log.Log("----rate limits: requests remaining: %d, used: %d, seconds left: %d", remaining, used, seconds)

	delay := poller.CalculatePollingRate(remaining, seconds)
	log.Log("----calulated polling delay to: %s", delay)

	return delay
}
