package constant

// SortType defines the available sorting options for posts.
type SortType string

const (
	SortTypeNew SortType = "new"
	SortTypeTop SortType = "top"
	SortTypeHot SortType = "hot"
)
