package book

import "github.com/shkuran/go-library/db"

func getBookById(id int64) (Book, error) {
	var book Book
	query := `
	SELECT * FROM books 
	WHERE id = $1
	`
	row := db.DB.QueryRow(query, id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.PublicationYear, &book.AvailableCopies)
	if err != nil {
		return book, err
	}

	return book, nil
}

func updateAvailableCopies(id, copies int64) error {
	query := `
	UPDATE books
	SET available_copies = $1
	WHERE id = $2
	`

	_, err := db.DB.Exec(query, copies, id)

	return err
}

func saveBook(book Book) error {
	query := `
	INSERT INTO books (title, author, isbn, publication_year, available_copies) 
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.DB.Exec(query, book.Title, book.Author, book.ISBN, book.PublicationYear, book.AvailableCopies)
	if err != nil {
		return err
	}

	return nil
}

func getBooks() ([]Book, error) {
	query := "SELECT * FROM books"
	rows, err := db.DB.Query(query)
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
