package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

// InitDB initializes the database by creating a new file, then opening a connection to it and creating a new user table if it doesn't already exists
func InitDB(dbFileName string) {
	CreateNewFile(dbFileName)
	DB = OpenConnection(dbFileName)
	CreateNewUserTable(DB)
	CreateNewResetPasswordTable(DB)
}

// CreateNewFile creates a new db file if doesn't already exists
func CreateNewFile(dbFileName string) {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		fmt.Println("Creating " + dbFileName)
		file, err := os.Create(dbFileName)
		if err != nil {
			log.Fatal("Error while creating the file")
		}
		file.Close()
		log.Println("Database file created")
	} else {
		log.Println("Using the already created database file")
	}
}

// CreateNewUserTable creates a new user table
func CreateNewUserTable(conn *sql.DB) {
	stmt, err := conn.Prepare(createUserTableSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	stmt.Exec()
}

// CreateNewResetPasswordTable creates a new reset password table
func CreateNewResetPasswordTable(conn *sql.DB) {
	stmt, err := conn.Prepare(createResetPasswordTableSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	stmt.Exec()
}

// OpenConnection opens a connection to the database
func OpenConnection(name string) *sql.DB {
	conn, err := sql.Open("sqlite3", name)

	if err != nil {
		log.Fatal(err.Error())
	}

	return conn
}
