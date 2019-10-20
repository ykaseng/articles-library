package database

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/ykaseng/articles-library/models"
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("database_dsn", fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", viper.GetString("DATABASE_URI_SCHEME"), viper.GetString("DATABASE_USER"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_HOST"), viper.GetInt("DATABASE_PORT"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_SSL")))
}

func TestNewArticleStore(t *testing.T) {}

func TestGet(t *testing.T) {
	tt := []struct {
		name     string
		seed     string
		id       int
		expected []models.Article
	}{

		{
			name:     "database no records",
			seed:     "",
			id:       1,
			expected: []models.Article(nil),
		},
		{
			name: "database has one record",
			seed: "WITH author AS (INSERT INTO authors(name) VALUES ('Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Test Title', 'Test Content', (SELECT author.id FROM author))",
			id:   1,
			expected: []models.Article{
				models.Article{
					ArticleID: models.ArticleID{
						ID: 1,
					},
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
			},
		},
		{
			name: "database has multiple records",
			seed: "WITH author AS (INSERT INTO authors(name) VALUES ('Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Test Title', 'Test Content', (SELECT author.id FROM author));WITH author AS (INSERT INTO authors(name) VALUES ('Another Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Another Test Title', 'Another Test Content', (SELECT author.id FROM author))",
			id:   2,
			expected: []models.Article{
				{
					ArticleID: models.ArticleID{
						ID: 2,
					},
					Title:   "Another Test Title",
					Content: "Another Test Content",
					Author:  "Another Test Author",
				},
			},
		},
		{
			name:     "record does not exist",
			seed:     "WITH author AS (INSERT INTO authors(name) VALUES ('Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Test Title', 'Test Content', (SELECT author.id FROM author));WITH author AS (INSERT INTO authors(name) VALUES ('Another Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Another Test Title', 'Another Test Content', (SELECT author.id FROM author))",
			id:       3,
			expected: []models.Article(nil),
		},
	}

	db, err := DBConn()
	if err != nil {
		t.Fatalf("open database connection: %v", err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := db.Begin()
			if err != nil {
				t.Errorf("failed to begin transaction: %v", err)
			}

			defer func() {
				tx.Rollback()
				restartSerial(t, db)
			}()

			if len(tc.seed) > 0 {
				if _, err := tx.Exec(tc.seed); err != nil {
					t.Errorf("failed to seed: %v", err)
				}
			}

			actual, err := (&ArticleStore{db: tx}).Get(tc.id)
			if err != nil {
				t.Errorf("getAll failed: %v", err)
			}

			assert.Equal(t, tc.expected, *actual)
		})
	}
}

func TestGetAll(t *testing.T) {
	tt := []struct {
		name     string
		seed     string
		expected []models.Article
	}{

		{
			name:     "database no records",
			seed:     "",
			expected: []models.Article(nil),
		},
		{
			name: "database has one record",
			seed: "WITH author AS (INSERT INTO authors(name) VALUES ('Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Test Title', 'Test Content', (SELECT author.id FROM author))",
			expected: []models.Article{
				models.Article{
					ArticleID: models.ArticleID{
						ID: 1,
					},
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
			},
		},
		{
			name: "database has multiple records",
			seed: "WITH author AS (INSERT INTO authors(name) VALUES ('Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Test Title', 'Test Content', (SELECT author.id FROM author));WITH author AS (INSERT INTO authors(name) VALUES ('Another Test Author') RETURNING id) INSERT INTO articles(title, content, author_id) VALUES('Another Test Title', 'Another Test Content', (SELECT author.id FROM author))",
			expected: []models.Article{
				{
					ArticleID: models.ArticleID{
						ID: 1,
					},
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
				{
					ArticleID: models.ArticleID{
						ID: 2,
					},
					Title:   "Another Test Title",
					Content: "Another Test Content",
					Author:  "Another Test Author",
				},
			},
		},
	}

	db, err := DBConn()
	if err != nil {
		t.Fatalf("open database connection: %v", err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := db.Begin()
			if err != nil {
				t.Errorf("failed to begin transaction: %v", err)
			}

			defer func() {
				tx.Rollback()
				restartSerial(t, db)
			}()

			if len(tc.seed) > 0 {
				if _, err := tx.Exec(tc.seed); err != nil {
					t.Errorf("failed to seed: %v", err)
				}
			}

			actual, err := (&ArticleStore{db: tx}).GetAll()
			if err != nil {
				t.Errorf("getAll failed: %v", err)
			}

			assert.Equal(t, tc.expected, *actual)
		})
	}
}

func TestPost(t *testing.T) {
	tt := []struct {
		name     string
		article  models.Article
		expected models.Article
	}{
		{
			name: "post record",
			article: models.Article{
				Title:   "Test Title",
				Author:  "Test Author",
				Content: "Test Content",
			},
			expected: models.Article{
				ArticleID: models.ArticleID{
					ID: 1,
				},
				Title:   "Test Title",
				Author:  "Test Author",
				Content: "Test Content",
			},
		},
	}

	db, err := DBConn()
	if err != nil {
		t.Fatalf("open database connection: %v", err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := db.Begin()
			if err != nil {
				t.Errorf("failed to begin transaction: %v", err)
			}

			defer func() {
				tx.Rollback()
				restartSerial(t, db)
			}()

			articleStore := &ArticleStore{db: tx}
			articleID, err := articleStore.Post(&tc.article)
			if err != nil {
				t.Errorf("post failed: %v", err)
			}

			actual, err := articleStore.Get(articleID.ID)
			if err != nil {
				t.Errorf("get failed: %v", err)
			}

			assert.Equal(t, []models.Article{tc.expected}, *actual)
		})
	}
}

func restartSerial(t *testing.T, db orm.DB) {
	_, err := db.Exec(`TRUNCATE articles RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Errorf("could not restart serial: %v", err)
	}
}
