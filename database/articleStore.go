package database

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/ykaseng/articles-library/models"
)

// ArticleStore implements database operations for article management.
type ArticleStore struct {
	db orm.DB
}

// NewArticleStore returns an ArticleStore.
func NewArticleStore(db orm.DB) *ArticleStore {
	return &ArticleStore{
		db: db,
	}
}

// Get an article by ID.
func (s *ArticleStore) Get(id int) (*[]models.Article, error) {
	q := `
	SELECT ar.id, ar.title, ar.content, au.name AS author FROM articles ar INNER JOIN authors au ON ar.author_id = au.id WHERE ar.id = ?
	`

	var a []models.Article
	if _, err := s.db.QueryOne(&a, q, id); err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		}
	}

	return &a, nil
}

// GetAll gets all articles.
func (s *ArticleStore) GetAll() (*[]models.Article, error) {
	q := `
	SELECT ar.id, ar.title, ar.content, au.name AS author FROM articles ar INNER JOIN authors au ON ar.author_id = au.id
	`

	var a []models.Article
	if _, err := s.db.Query(&a, q); err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		}
	}

	return &a, nil
}

// Post inserts an article into the database and returns the last insert id.
func (s *ArticleStore) Post(article *models.Article) (*models.ArticleID, error) {
	q := `
		WITH author AS (INSERT INTO authors(name) VALUES (?) RETURNING id) INSERT INTO articles(title, content, author_id) VALUES(?, ?, (SELECT author.id FROM author)) RETURNING id
	`

	var articleID models.ArticleID
	if _, err := s.db.QueryOne(&articleID.ID, q, article.Author, article.Title, article.Content); err != nil {
		return nil, err
	}

	return &articleID, nil
}
