package book

import (
	"database/sql"
)

type Repository interface {
	GetById(id int64) (Book, error)
	UpdateAvailableCopies(id, copies int64) error
	Save(book Book) error
	GetAll() ([]Book, error)
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetById(id int64) (Book, error) {
	var book Book
	query := `
	SELECT * FROM books 
	WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (r *Repo) UpdateAvailableCopies(id, copies int64) error {
	query := `
	UPDATE books
	SET available_copies = $1
	WHERE id = $2
	`

	_, err := r.db.Exec(query, copies, id)

	return err
}

func (r *Repo) Save(book Book) error {
	query := `
	INSERT INTO books (title, author, isbn, publication_year, available_copies) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, book.Title, book.Author, book.ISBN, book.PublicationYear, book.AvailableCopies)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetAll() ([]Book, error) {
	query := "SELECT * FROM books"
	rows, err := r.db.Query(query)
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
