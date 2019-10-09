// Package models contains application specific entities.
package models

import (
	"github.com/go-ozzo/ozzo-validation"
)

// Article holds specific application settings linked to an Account.
type Article struct {
	ID      int    `json:"-"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Validate validates Article struct and returns validation errors.
func (a *Article) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Content, validation.Required),
		validation.Field(&a.Author, validation.Required),
	)
}
