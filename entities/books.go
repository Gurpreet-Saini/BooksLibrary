package entities

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        Author `json:"author,omitempty"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"published_date"`
}

type ContextKey string

const (
	Title         ContextKey = "title"
	IncludeAuthor ContextKey = "includeAuthor"
	Id            ContextKey = "id"
	FirstName     ContextKey = "FirstName"
)
