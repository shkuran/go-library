package book

import (
	"database/sql"
)

type Repository interface {
	GetBookById(id int64) (Book, error)
	UpdateAvailableCopies(id, copies int64) error
	SaveBook(book Book) error
	GetBooks() ([]Book, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (svc *MySQLRepository) GetBookById(id int64) (Book, error) {
	var book Book

	row := svc.db.QueryRow("SELECT * FROM books WHERE id = ?", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (svc *MySQLRepository) UpdateAvailableCopies(id, copies int64) error {
	query := `
	UPDATE books
	SET available_copies = ?
	WHERE id = ?
	`

	_, err := svc.db.Exec(query, copies, id)

	return err
}

func (svc *MySQLRepository) SaveBook(book Book) error {
	query := `
	INSERT INTO books (title, author, isbn, publication_year, available_copies) 
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := svc.db.Exec(query, book.Title, book.Author, book.ISBN, book.PublicationYear, book.AvailableCopies)
	if err != nil {
		return err
	}

	return nil
}

func (svc *MySQLRepository) GetBooks() ([]Book, error) {
	rows, err := svc.db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
