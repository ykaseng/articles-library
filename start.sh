#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS authors (id SERIAL, name VARCHAR(255), PRIMARY KEY(id));
    CREATE TABLE IF NOT EXISTS articles (id SERIAL, title TEXT, content TEXT, author_id INT, PRIMARY KEY(id), FOREIGN KEY(author_id) REFERENCES authors(id));
EOSQL