package poller

import (
	"encoding/json"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Post struct {
	Subreddit   string
	Title       string
	Author      string
	Score       int
	UpvoteRatio float32
	Permalink   string
}

func MakePost(p *reddit.Post) Post {
	return Post{
		Subreddit:   p.SubredditName,
		Title:       p.Title,
		Author:      p.Author,
		Score:       p.Score,
		UpvoteRatio: p.UpvoteRatio,
		Permalink:   p.Permalink,
	}
}

type Posts []Post

func MakePosts(redditPosts ...*reddit.Post) []Post {
	posts := make([]Post, len(redditPosts))
	for i := 0; i < len(posts); i++ {
		posts[i] = MakePost(redditPosts[i])
	}
	return Posts(posts)
}

func (p *Posts) MarshalJSON() ([]byte, error) {
	type Alias Posts

	bytes, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})

	if err != nil {
		return nil, fmt.Errorf("could not marshal posts: %s", err)
	}
	return bytes, nil
}
