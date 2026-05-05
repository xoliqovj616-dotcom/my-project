package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS todos(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		work TEXT NOT NULL,
		time TEXT,
		completed INTEGER DEFAULT 0
	);`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
