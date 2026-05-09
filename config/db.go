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
		time TEXT NOT NULL,
		completed INTEGER DEFAULT 0,
		user_id INTEGER NOT NULL
	);`
	queryUser := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);`
	_, err = DB.Exec(queryUser)
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
