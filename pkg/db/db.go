package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// scheduler table and scheduler_date index creation script
const schema = `
	CREATE TABLE scheduler ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT DEFAULT "",
		repeat VARCHAR(128) DEFAULT ""
	);
	CREATE INDEX scheduler_date ON scheduler (date);
	`

var DB *sql.DB

// Database initialization function
func Init(dbFile string) error {

	var install bool

	// Validation if database file exists in a root directory
	_, err := os.Stat(dbFile)

	if err != nil {
		install = true
	}

	// Opening a database
	DB, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return err
	}

	if install == true {

		// Creating schema for a new database
		_, err = DB.Exec(schema)

		if err != nil {
			return err
		}
	}


	return nil

}
