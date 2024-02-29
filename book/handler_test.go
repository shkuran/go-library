package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func InitBookTable() []Book {

	books := []Book{
		{
			ID:              1,
			Title:           "The Great Gatsby",
			Author:          "F. Scott Fitzgerald",
			ISBN:            "978-0743273565",
			PublicationYear: 1925,
			AvailableCopies: 10,
		},
		{
			ID:              2,
			Title:           "To Kill a Mockingbird",
			Author:          "Harper Lee",
			ISBN:            "978-0061120084",
			PublicationYear: 1960,
			AvailableCopies: 8,
		},
	}
	return books
}

type MockBookRepo struct {
	books []Book
}

func NewMockBookRepo(books []Book) *MockBookRepo {
	return &MockBookRepo{books: books}
}

func (r *MockBookRepo) getAll() ([]Book, error) {
	return r.books, nil
	// return nil, errors.New("Could not fetch books!")
}

func (r *MockBookRepo) GetById(id int64) (Book, error) {
	return r.books[id], nil
}

func (r *MockBookRepo) UpdateAvailableCopies(id, copies int64) error {
	r.books[id].AvailableCopies = copies
	return nil
}

func (r *MockBookRepo) save(book Book) error {
	r.books = append(r.books, book)
	return nil
}

func TestGetBooks(t *testing.T) {
	books := InitBookTable()
	book_repo := NewMockBookRepo(books)
	book_handler := NewHandler(book_repo)

	router := gin.Default()

	router.GET("/books", book_handler.GetBooks)

	// Perform a test request
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	var responseBooks []Book
	err = json.Unmarshal(w.Body.Bytes(), &responseBooks)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(books, responseBooks) {
		t.Errorf("Expected %+v; got %+v", books, responseBooks)
	}

}
