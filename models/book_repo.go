package models

import "github.com/shkuran/go-library/db"

type Book struct {
	Id              int64  `json:"id" db:"id"`
	Title           string `json:"title" db:"title" binding:"required"`
	Author          string `json:"author" db:"author" binding:"required"`
	ISBN            string `json:"isbn" db:"isbn" binding:"required"`
	PublicationYear int64  `json:"publication_year" db:"publication_year" binding:"required"`
	AvailableCopies int64  `json:"available_copies" db:"available_copies" binding:"required"`
}

func GetBooks() ([]Book, error) {
	rows, err := db.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func GetBookById(id int64) (Book, error) {
	var book Book

	row := db.DB.QueryRow("SELECT * FROM books WHERE id = ?", id)
	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
	if err != nil {
		return book, err
	}

	return book, nil
}

func SaveBook(book Book) error {
	query := `
	INSERT INTO books (title, author, isbn, publication_year, available_copies) 
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.DB.Exec(query, book.Title, book.Author, book.ISBN, book.PublicationYear, book.AvailableCopies)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAvailableCopies(id, copies int64) error {
	query := `
	UPDATE books
	SET available_copies = ?
	WHERE id = ?
	`

	_, err := db.DB.Exec(query, copies, id)

	return err
}
