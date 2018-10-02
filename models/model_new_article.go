package models

type NewArticle struct {

	Title string `json:"title"`

	Description string `json:"description"`

	Body string `json:"body"`

	TagList []string `json:"tagList,omitempty"`
}
