package poller

import "sort"

type Authors struct {
	AuthorsMap map[AuthorID]*Author
}

func MakeAuthors() Authors {
	m := make(map[AuthorID]*Author)
	return Authors{m}
}

func (a Authors) CountPosts(posts Posts) {
	for _, post := range posts {

		ID := post.Author.ID

		author, exists := a.AuthorsMap[ID]
		if !exists {

			a.AuthorsMap[ID] = post.Author
			author, _ = a.AuthorsMap[ID]
		}
		author.TotalPosts += 1
	}
}

func (a Authors) TopTenAuthors() []*Author {
	return a.TopAuthorsForCount(10)
}

func (a Authors) TopAuthorsForCount(count int) []*Author {
	if count > len(a.AuthorsMap) {
		count = len(a.AuthorsMap)
	}
	authors := make([]*Author, 0, len(a.AuthorsMap))

	for _, author := range a.AuthorsMap {
		authors = append(authors, author)
	}

	sort.Slice(authors, func(i, j int) bool {
		x, y := authors[i], authors[j]

		if x.TotalPosts == y.TotalPosts {
			return x.ID < y.ID
		}

		return x.TotalPosts > y.TotalPosts
	})

	return authors[:count]
}
