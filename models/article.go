// Package models contains application specific entities.
package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Article holds specific application settings linked to an Article.
type Article struct {
	ArticleID
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Validate validates Article struct and returns validation errors.
func (a *Article) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Content, validation.Required),
		validation.Field(&a.Author, validation.Required, validation.Length(5, 255)),
	)
}

// ArticleID holds specific application settings linked to an ArticleID.
type ArticleID struct {
	ID int `json:"id"`
}
