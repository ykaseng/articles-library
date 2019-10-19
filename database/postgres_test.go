package database

import (
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.AutomaticEnv()
}

func TestDBConn(t *testing.T) {
	tt := []struct {
		name     string
		dsn      string
		expected struct {
			err error
		}
	}{
		{
			name: "bad user",
			dsn:  fmt.Sprintf("postgres://baduser:badpassword@%s:5432/library?sslmode=disable", viper.GetString("DATABASE_HOST")),
			expected: struct {
				err error
			}{
				err: errors.New("FATAL #28P01 password authentication failed for user \"baduser\""),
			},
		},
		{
			name: "bad password",
			dsn:  fmt.Sprintf("postgres://postgres:badpassword@%s:5432/library?sslmode=disable", viper.GetString("DATABASE_HOST")),
			expected: struct {
				err error
			}{
				err: errors.New("FATAL #28P01 password authentication failed for user \"postgres\""),
			},
		},
		{
			name: "correct dsn",
			dsn:  viper.GetString("database_dsn"),
			expected: struct {
				err error
			}{
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer resetEnv()
			viper.SetDefault("database_dsn", tc.dsn)

			_, err := DBConn()
			if err != nil {
				if tc.expected.err != nil {
					assert.Equal(t, tc.expected.err.Error(), err.Error())
					return
				}

				t.Errorf("connection failed: %v", err)
			}
		})
	}
}

func resetEnv() {
	viper.SetDefault("database_dsn", fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", viper.GetString("DATABASE_URI_SCHEME"), viper.GetString("DATABASE_USER"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_HOST"), viper.GetInt("DATABASE_PORT"), viper.GetString("DATABASE_NAME"), viper.GetString("DATABASE_SSL")))
}
