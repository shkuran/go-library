package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "root",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "library",
		ParseTime: true,
	}

	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")

	createTables()
}

func createTables() {
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
	_, err := DB.Exec(createBooksTable)
	if err != nil {
		panic("Cannot create books table!")
	}
	log.Println("Table books was successfuly created!")

	createClientsTable := `
	CREATE TABLE IF NOT EXISTS clients (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE
	);
	`
	_, err = DB.Exec(createClientsTable)
	if err != nil {
		panic("Cannot create clients table!")
	}
	log.Println("Table clients was successfuly created!")

	createReservationsTable := `
	CREATE TABLE IF NOT EXISTS reservations (
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		book_id INT NOT NULL,
		client_id INT NOT NULL,
		checkout_date DATETIME NOT NULL,
		return_date DATETIME,
		FOREIGN KEY (book_id) REFERENCES books(id),
    	FOREIGN KEY (client_id) REFERENCES clients(id)
	);
	`
	_, err = DB.Exec(createReservationsTable)
	if err != nil {
		panic("Cannot create reservations table!")
	}
	log.Println("Table reservations was successfuly created!")
}
