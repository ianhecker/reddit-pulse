package poller

import (
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Post struct {
	Subreddit   string
	Title       string
	Author      Author
	Score       int
	UpvoteRatio float32
	Permalink   string
}

func MakePost(p *reddit.Post) Post {
	author := Author{Name: p.Author, ID: p.AuthorID}

	return Post{
		Subreddit:   p.SubredditName,
		Title:       p.Title,
		Author:      author,
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
