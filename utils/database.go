package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToSqlite() (*sql.DB, error) {
	dbName := "local.db"
	return sql.Open("sqlite3", dbName)
}
