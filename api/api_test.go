package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ykaseng/articles-library/api/app"
	"github.com/ykaseng/articles-library/models"
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("database_dsn", fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", viper.GetString("DATABASE_URI_SCHEME"), viper.GetString("DATABASE_USER"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_HOST"), viper.GetInt("DATABASE_PORT"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_SSL")))
}

func TestRouter(t *testing.T) {
	type getAllArticlesResponse struct {
		app.Status
		Data []models.Article `json:"data"`
	}

	tt := []struct {
		name     string
		method   string
		endpoint string
		expected getAllArticlesResponse
	}{
		{
			name:     "correct endpoint",
			method:   "GET",
			endpoint: "/articles",
			expected: getAllArticlesResponse{
				Status: app.Status{
					Code:    http.StatusOK,
					Message: "SUCCESS",
				},
				Data: []models.Article(nil),
			},
		},
		{
			name:     "invalid endpoint",
			method:   "GET",
			endpoint: "/artichokes",
			expected: getAllArticlesResponse{
				Status: app.Status{
					Code:    http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
				},
				Data: []models.Article(nil),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			api, err := New()
			if err != nil {
				t.Errorf("failed to create api : %v", err)
			}

			srv := httptest.NewServer(api)
			defer srv.Close()

			res := testRequest(t, srv, tc.method, tc.endpoint, nil)
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("read response failed: %v", err)
			}

			var actual getAllArticlesResponse
			if err := json.Unmarshal(b, &actual); err != nil {
				t.Errorf("unmarshal response failed: %v", err)
			}

			assert.Equal(t, tc.expected, actual)
		})
	}

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return resp
}
