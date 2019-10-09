package database

import (
	"fmt"

	"github.com/go-pg/pg"

	"github.com/ykaseng/articles-library/models"
)

// ArticleStore implements database operations for account management by user.
type ArticleStore struct {
	db *pg.DB
}

// NewArticleStore returns an ArticleStore.
func NewArticleStore(db *pg.DB) *ArticleStore {
	return &ArticleStore{
		db: db,
	}
}

// Get an account by ID.
func (s *ArticleStore) Get(id int) (*models.Article, error) {
	return &models.Article{}, fmt.Errorf("not implemented")
}

// Get an account by ID.
func (s *ArticleStore) GetAll() (*[]models.Article, error) {
	return &[]models.Article{}, fmt.Errorf("not implemented")
}

// Get an account by ID.
func (s *ArticleStore) Post(article *models.Article) (int, error) {
	return 99, nil
	// return 0, fmt.Errorf("not implemented")
}
