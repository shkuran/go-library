package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/book"
)

func TestGetBooks(t *testing.T) {

	testCases := []struct {
		testName         string
		booksInDB        []book.Book
		expectedCode     int
		expectedBooks    []book.Book
		expectedErrorMsg string
	}{
		// Case 1: GetBooks returns []Book
		{
			testName:         "Success",
			booksInDB:        []book.Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}},
			expectedCode:     http.StatusOK,
			expectedBooks:    []book.Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}},
			expectedErrorMsg: "",
		},
		// Case 2: GetBooks returns an error
		{
			testName:         "Error",
			booksInDB:        []book.Book{},
			expectedCode:     http.StatusInternalServerError,
			expectedBooks:    nil,
			expectedErrorMsg: "Could not fetch books!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			book_repo := NewMockBookRepo(tc.booksInDB)
			book_handler := book.NewHandler(book_repo)

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

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
			}

			if tc.expectedErrorMsg != "" {
				// Check if the response contains the expected error message
				var response map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}
			} else {
				// Check if the response matches the expected books
				var responseBooks []book.Book
				err = json.Unmarshal(w.Body.Bytes(), &responseBooks)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.expectedBooks, responseBooks) {
					t.Errorf("Expected %+v; got %+v", tc.expectedBooks, responseBooks)
				}
			}
		})
	}

}
