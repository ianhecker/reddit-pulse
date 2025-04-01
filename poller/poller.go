package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Poll struct {
	Posts    Posts
	Response *Response
	Error    error
}

type Poller struct {
	Client             *reddit.Client
	AvailableRequests  int
	TimeFrameInSeconds int
	PollingInterval    time.Duration
}

type PollFunc func(ctx context.Context) Poll

func NewPoller(
	credentials reddit.Credentials,
	opt ...reddit.Opt,
) (*Poller, error) {
	client, err := reddit.NewClient(credentials)
	if err != nil {
		return nil, fmt.Errorf("could not make reddit client: %s", err)
	}
	return NewPollerFromRaw(client, 60, 60, time.Duration(1*time.Second)), nil
}

func NewPollerFromRaw(
	client *reddit.Client,
	availableRequests int,
	timeFrameInSeconds int,
	pollingInterval time.Duration,
) *Poller {
	return &Poller{
		Client:             client,
		AvailableRequests:  availableRequests,
		TimeFrameInSeconds: timeFrameInSeconds,
		PollingInterval:    pollingInterval,
	}
}

func (p Poller) TopPosts(
	ctx context.Context,
	subreddit string,
	limit int,
) Poll {
	options := &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: limit,
		},
		Time: "all",
	}
	posts, response, err := p.Client.Subreddit.TopPosts(ctx, subreddit, options)
	if err != nil {
		return Poll{
			Posts:    nil,
			Response: NewResponse(response),
			Error:    fmt.Errorf("could not fetch top posts: %s", err),
		}
	}
	return Poll{
		Posts:    MakePosts(posts...),
		Response: NewResponse(response),
		Error:    nil,
	}
}
