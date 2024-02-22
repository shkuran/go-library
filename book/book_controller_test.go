package book

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetBooks(t *testing.T) {
	testBooks := []Book{
		{
			Id:              1,
			Title:           "Test Book 1",
			Author:          "Test Author 1",
			ISBN:            "isbn",
			PublicationYear: 1986,
			AvailableCopies: 5,
		},
		{
			Id:              2,
			Title:           "Test Book 2",
			Author:          "Test Author 2",
			ISBN:            "isbn",
			PublicationYear: 1986,
			AvailableCopies: 5,
		},
	}

	// Save the original getBooks function
	originalGetBooksFunc := getBooksFunc

	// Mock the getBooks function
	getBooksFunc = func() ([]Book, error) {
		return testBooks, nil
	}
	defer func() {
		// Restore the original getBooksFunc after the test
		getBooksFunc = originalGetBooksFunc
	}()

	// Create a new Gin router
	router := gin.Default()

	// Set up the test route using the GetBooks handler
	router.GET("/books", GetBooks)

	// Perform a GET request to the /books endpoint
	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}

	// Parse the response body to check its content
	var responseBooks []Book
	err = json.Unmarshal(w.Body.Bytes(), &responseBooks)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the response matches the expected mock data
	expectedBooks := testBooks
	if !booksSliceEqual(expectedBooks, responseBooks) {
		t.Errorf("Expected %+v; got %+v", expectedBooks, responseBooks)
	}
}

// booksSliceEqual checks if two slices of books are equal
func booksSliceEqual(slice1, slice2 []Book) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

// func TestAddBook(t *testing.T) {
// 	r := setupRouter()

// 	newBook := Book{
// 		Title: "Test Book",
// 		Author: "Test Author",
// 		ISBN: "isbn",
// 		PublicationYear: 1986,
// 		AvailableCopies: 5,
// 	}
// 	newBookJSON, _ := json.Marshal(newBook)
// 	reqAdd := httptest.NewRequest("POST", "/books", bytes.NewBuffer(newBookJSON))
// 	reqAdd.Header.Set("Content-Type", "application/json")
// 	wAdd := httptest.NewRecorder()
// 	r.ServeHTTP(wAdd, reqAdd)

// 	if wAdd.Code != http.StatusCreated {
// 		t.Errorf("Expected status %d; got %d", http.StatusCreated, wAdd.Code)
// 	}
// }

// func setupRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.GET("/books", GetBooks)
// 	r.POST("/books", AddBook)
// 	return r
// }
