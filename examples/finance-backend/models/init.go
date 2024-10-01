// init.go
package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase initializes the SQLite database and creates the necessary tables
func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./finance_Db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	createTables := `
	CREATE TABLE IF NOT EXISTS income (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		amount REAL NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS expense (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		amount REAL NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`
	_, err = DB.Exec(createTables)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
	}
}
