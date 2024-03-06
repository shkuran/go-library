package book

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetById(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// Mock database expectations
	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "publication_year", "available_copies"}).
		AddRow(1, "Test Book", "Test Author", "123456789", 2022, 5)

	mock.ExpectQuery("SELECT \\* FROM books WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	// Perform the actual function call
	result, err := repo.GetById(1)

	// Check the result and assertions
	assert.NoError(t, err)
	assert.Equal(t, Book{
		ID:              1,
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "123456789",
		PublicationYear: 2022,
		AvailableCopies: 5,
	}, result)

	mock.ExpectQuery("SELECT \\* FROM books WHERE id = \\$1").
		WithArgs(1).
		WillReturnError(errors.New("simulate err from db"))

	// Perform the actual function call
	result, err = repo.GetById(1)

	// Check the result and assertions
	assert.Error(t, err)
	assert.Equal(t, Book{}, result)
}

func TestUpdateAvailableCopies(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// Mock database expectations
	mock.ExpectExec("UPDATE books SET available_copies = \\$1 WHERE id = \\$2").
		WithArgs(10, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Perform the actual function call
	err = repo.UpdateAvailableCopies(1, 10)

	// Check the result and assertions
	assert.NoError(t, err)
}

func TestSave(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// Mock database expectations
	mock.ExpectExec("INSERT INTO books").
		WithArgs("New Book", "New Author", "987654321", 2023, 5).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Perform the actual function call
	err = repo.Save(Book{
		Title:           "New Book",
		Author:          "New Author",
		ISBN:            "987654321",
		PublicationYear: 2023,
		AvailableCopies: 5,
	})

	// Check the result and assertions
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// Mock database expectations
	rows := sqlmock.NewRows([]string{"id", "title", "author", "isbn", "publication_year", "available_copies"}).
		AddRow(1, "Book 1", "Author 1", "111", 2022, 3).
		AddRow(2, "Book 2", "Author 2", "222", 2023, 7)

	mock.ExpectQuery("SELECT \\* FROM books").
		WillReturnRows(rows)

	// Perform the actual function call
	result, err := repo.GetAll()

	// Check the result and assertions
	assert.NoError(t, err)
	assert.Equal(t, []Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "111", PublicationYear: 2022, AvailableCopies: 3},
		{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "222", PublicationYear: 2023, AvailableCopies: 7},
	}, result)
}
