package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ykaseng/articles-library/database"
	"github.com/ykaseng/articles-library/models"
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("database_dsn", fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", viper.GetString("DATABASE_URI_SCHEME"), viper.GetString("DATABASE_USER"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_HOST"), viper.GetInt("DATABASE_PORT"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_SSL")))
}

func TestGetAll(t *testing.T) {
	type getAllArticlesResponse struct {
		Status
		Data []models.Article `json:"data"`
	}

	tt := []struct {
		name     string
		seeds    []models.Article
		expected getAllArticlesResponse
	}{
		{
			name:  "database no records",
			seeds: []models.Article{},
			expected: getAllArticlesResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article(nil),
			},
		},
		{
			name: "database one record",
			seeds: []models.Article{
				{
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
			},
			expected: getAllArticlesResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article{
					{
						ArticleID: models.ArticleID{
							ID: 1,
						},
						Title:   "Test Title",
						Content: "Test Content",
						Author:  "Test Author",
					},
				},
			},
		},
		{
			name: "database multiple record",
			seeds: []models.Article{
				{
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
				{
					Title:   "Another Test Title",
					Content: "Another Test Content",
					Author:  "Another Test Author",
				},
			},
			expected: getAllArticlesResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article{
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
		},
	}

	db, err := database.DBConn()
	if err != nil {
		t.Fatalf("open database connection: %v", err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := db.Begin()
			if err != nil {
				t.Errorf("failed to begin transaction : %v", err)
			}

			defer func() {
				tx.Rollback()
				restartSerial(t, db)
			}()

			article := NewArticleResource(database.NewArticleStore(tx))
			for _, seed := range tc.seeds {
				_, err := article.Store.Post(&seed)
				if err != nil {
					t.Errorf("failed to seed: %v", err)
				}
			}

			req, err := http.NewRequest("GET", "localhost:8080/api/v1/articles", nil)
			if err != nil {
				t.Errorf("create request failed: %v", err)
			}

			rec := httptest.NewRecorder()
			article.getAll(rec, req)

			res := rec.Result()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("read response failed: %v", err)
			}

			var actual getAllArticlesResponse
			if err := json.Unmarshal(b, &actual); err != nil {
				t.Errorf("unmarshal response failed: %v", err)
			}

			assert.Equal(t, tc.expected.Code, actual.Code)
			assert.Equal(t, tc.expected.Message, actual.Message)
			assert.Equal(t, tc.expected.Data, actual.Data)
		})
	}
}

func TestGet(t *testing.T) {
	type getArticleResponse struct {
		Status
		Data []models.Article `json:"data"`
	}

	tt := []struct {
		name     string
		seeds    []models.Article
		id       int
		expected getArticleResponse
	}{
		{
			name:  "database no records",
			seeds: []models.Article{},
			expected: getArticleResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article(nil),
			},
		},
		{
			name: "database one record",
			seeds: []models.Article{
				{
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
			},
			id: 1,
			expected: getArticleResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article{
					{
						ArticleID: models.ArticleID{
							ID: 1,
						},
						Title:   "Test Title",
						Content: "Test Content",
						Author:  "Test Author",
					},
				},
			},
		},
		{
			name: "database multiple record",
			seeds: []models.Article{
				{
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
				{
					Title:   "Another Test Title",
					Content: "Another Test Content",
					Author:  "Another Test Author",
				},
			},
			id: 2,
			expected: getArticleResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article{
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
		},
		{
			name: "record does not exist",
			seeds: []models.Article{
				{
					Title:   "Test Title",
					Content: "Test Content",
					Author:  "Test Author",
				},
				{
					Title:   "Another Test Title",
					Content: "Another Test Content",
					Author:  "Another Test Author",
				},
			},
			id: 3,
			expected: getArticleResponse{
				Status: Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article(nil),
			},
		},
	}

	db, err := database.DBConn()
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

			article := NewArticleResource(database.NewArticleStore(tx))
			for _, seed := range tc.seeds {
				_, err := article.Store.Post(&seed)
				if err != nil {
					t.Errorf("failed to seed: %v", err)
				}
			}

			req, err := http.NewRequest("GET", fmt.Sprintf("localhost:8080/api/v1/articles/%d", tc.id), nil)
			if err != nil {
				t.Errorf("create request failed: %v", err)
			}

			rec := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("articleID", fmt.Sprintf("%d", tc.id))

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			article.get(rec, req)
			res := rec.Result()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("read response failed: %v", err)
			}

			var actual getArticleResponse
			if err := json.Unmarshal(b, &actual); err != nil {
				t.Errorf("unmarshal response failed: %v", err)
			}

			assert.Equal(t, tc.expected.Code, actual.Code)
			assert.Equal(t, tc.expected.Message, actual.Message)
			assert.Equal(t, tc.expected.Data, actual.Data)
		})
	}
}

func restartSerial(t *testing.T, db orm.DB) {
	_, err := db.Exec(`TRUNCATE articles RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Errorf("could not restart serial: %v", err)
	}
}
