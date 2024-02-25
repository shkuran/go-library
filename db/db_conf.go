package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "root",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "library",
		ParseTime: true,
	}
	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, pingErr

	}
	log.Println("Connected!")

	return db, nil
}

func CreateTables(db *sql.DB) {
	createBooksTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		isbn VARCHAR(20),
		publication_year INT,
    	available_copies INT
	);
	`
	_, err := db.Exec(createBooksTable)
	if err != nil {
		panic("Cannot create books table!")
	}
	log.Println("Table books was successfuly created!")

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		panic("Cannot create users table!")
	}
	log.Println("Table users was successfuly created!")

	createReservationsTable := `
	CREATE TABLE IF NOT EXISTS reservations (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		book_id INT NOT NULL,
		user_id INT NOT NULL,
		checkout_date DATETIME NOT NULL,
		return_date DATETIME,
		FOREIGN KEY (book_id) REFERENCES books(id),
    	FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err = db.Exec(createReservationsTable)
	if err != nil {
		panic("Cannot create reservations table!")
	}
	log.Println("Table reservations was successfuly created!")
}
