package poller

type AuthorID string

type Author struct {
	Name       string
	ID         AuthorID
	TotalPosts *int
}

func MakeAuthor(name string, ID string) Author {
	total := 0
	return Author{Name: name, ID: AuthorID(ID), TotalPosts: &total}
}
