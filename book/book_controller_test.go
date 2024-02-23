package book

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetBooks(t *testing.T) {

	testCases := []struct {
		name             string
		getBooksFunc     func() ([]Book, error)
		expectedCode     int
		expectedBooks    []Book
		expectedErrorMsg string
	}{
		// Case 1: getBooks returns []Book
		{
			name: "Success",
			getBooksFunc: func() ([]Book, error) {
				return []Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}}, nil
			},
			expectedCode:     http.StatusOK,
			expectedBooks:    []Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}},
			expectedErrorMsg: "",
		},
		// Case 2: getBooks returns an error
		{
			name: "Error",
			getBooksFunc: func() ([]Book, error) {
				return nil, errors.New("Simulated error fetching books")
			},
			expectedCode:     http.StatusInternalServerError,
			expectedBooks:    nil,
			expectedErrorMsg: "Could not fetch books!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()

			router.GET("/books", func(context *gin.Context) {
				GetBooks(context, tc.getBooksFunc)
			})

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
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
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
				var responseBooks []Book
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
