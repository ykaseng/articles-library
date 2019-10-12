package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const authorsTable = `
CREATE TABLE authors (
	id SERIAL,
	name VARCHAR(255),
	
	PRIMARY KEY(id)
)`

const articlesTable = `
CREATE TABLE articles (
	id SERIAL,
	title TEXT,
	content TEXT,
	author_id INT,
	
	PRIMARY KEY(id),
	FOREIGN KEY(author_id) REFERENCES authors(id)
  )`

func init() {
	up := []string{
		authorsTable,
		articlesTable,
	}

	down := []string{
		`DROP TABLE articles`,
		`DROP TABLE authors`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
