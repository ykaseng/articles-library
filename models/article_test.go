package models

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tt := []struct {
		name    string
		article *Article
		err     string
	}{
		{"article missing no fields", &Article{Title: "TestTitle", Content: "TestContent", Author: "TestAuthor"}, ""},
		{"article missing title field", &Article{Content: "TestContent", Author: "TestAuthor"}, "title: cannot be blank."},
		{"article missing content field", &Article{Title: "TestTitle", Author: "TestAuthor"}, "content: cannot be blank."},
		{"article missing author field", &Article{Title: "TestTitle", Content: "TestContent"}, "author: cannot be blank."},
		{"article missing multiple fields", &Article{Title: "TestTitle"}, "author: cannot be blank; content: cannot be blank."},
		{"article missing all fields", &Article{}, "author: cannot be blank; content: cannot be blank; title: cannot be blank."},
		{"article invalid author field", &Article{Title: "TestTitle", Content: "TestContent", Author: "Test"}, "author: the length must be between 1 and 255."},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.article.Validate(); err != nil {
				if strings.Compare(tc.err, err.Error()) != 0 {
					t.Errorf("validate of %v should be %v; got %v", tc.name, tc.err, err)
				}
			}
		})
	}
}
