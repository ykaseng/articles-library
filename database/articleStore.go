package database

import (
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

// Get an article by ID.
func (s *ArticleStore) Get(id int) (*[]models.Article, error) {
	q := `
	SELECT ar.title, ar.content, au.name FROM articles ar INNER JOIN authors au ON ar.author_id = au.id WHERE ar.id = $1
	`

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var a models.Article
	if _, err = stmt.QueryOne(pg.Scan(&a.Title, &a.Content, &a.Author), id); err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		}

		return nil, nil
	}

	return &[]models.Article{a}, nil
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

		return nil, nil
	}

	return &a, nil
}

// Post inserts an article into the database and returns the last insert id.
func (s *ArticleStore) Post(article *models.Article) (*models.ArticleID, error) {
	q := `
		WITH author AS (INSERT INTO authors(name) VALUES ($1) RETURNING id) INSERT INTO articles(title, content, author_id) VALUES($2, $3, (SELECT author.id FROM author)) RETURNING id
	`

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var articleID models.ArticleID
	if _, err := stmt.QueryOne(pg.Scan(&articleID.ID), article.Author, article.Title, article.Content); err != nil {
		return nil, err
	}

	return &articleID, nil
}
