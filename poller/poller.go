package poller

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Poller struct {
	Client             *reddit.Client
	AvailableRequests  int
	TimeFrameInSeconds int
}

func NewPoller(
	credentials reddit.Credentials,
	opt ...reddit.Opt,
) (*Poller, error) {
	client, err := reddit.NewClient(credentials)
	if err != nil {
		return nil, fmt.Errorf("could not make reddit client: %s", err)
	}
	return NewPollerFromRaw(client, 1000, 60), nil
}

func NewPollerFromRaw(
	client *reddit.Client,
	availableRequests int,
	timeFrameInSeconds int,
) *Poller {
	return &Poller{
		Client:             client,
		AvailableRequests:  availableRequests,
		TimeFrameInSeconds: timeFrameInSeconds,
	}
}

func (p Poller) TopPosts(
	ctx context.Context,
	subreddit string,
	limit int,
) (Posts, *Response, error) {
	options := &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: limit,
		},
		Time: "all",
	}
	posts, response, err := p.Client.Subreddit.TopPosts(ctx, "golang", options)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch top posts: %s", err)
	}

	postsWrapper := Posts(posts)
	reponseWrapper := Response(*response)

	return postsWrapper, &reponseWrapper, nil
}
