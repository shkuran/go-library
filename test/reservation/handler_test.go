package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/book"
	"github.com/shkuran/go-library/reservation"
	test "github.com/shkuran/go-library/test/book"
)

func TestGetReservations(t *testing.T) {

	testCases := []struct {
		testName             string
		booksInDB            []book.Book
		reservationsInDB     []reservation.Reservation
		expectedCode         int
		expectedReservations []reservation.Reservation
		expectedErrorMsg     string
	}{
		// Case 1: GetReservation returns []Reservation
		{
			testName:             "Success",
			booksInDB:            []book.Book{},
			reservationsInDB:     []reservation.Reservation{{ID: 1, BookId: 1, UserId: 1}, {ID: 2, BookId: 2, UserId: 2}},
			expectedCode:         http.StatusOK,
			expectedReservations: []reservation.Reservation{{ID: 1, BookId: 1, UserId: 1}, {ID: 2, BookId: 2, UserId: 2}},
			expectedErrorMsg:     "",
		},
		// Case 2: GetReservation returns an error
		{
			testName:             "Error",
			booksInDB:            []book.Book{},
			reservationsInDB:     []reservation.Reservation{},
			expectedCode:         http.StatusInternalServerError,
			expectedReservations: nil,
			expectedErrorMsg:     "Could not fetch reservations!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			book_repo := test.NewMockBookRepo(tc.booksInDB)
			res_repo := NewMockReservationRepo(tc.reservationsInDB)
			res_handler := reservation.NewHandler(res_repo, book_repo)

			router := gin.Default()

			router.GET("/reservations", res_handler.GetReservations)

			// Perform a test request
			req, err := http.NewRequest("GET", "/reservations", nil)
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
				var responseReservations []reservation.Reservation
				err = json.Unmarshal(w.Body.Bytes(), &responseReservations)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tc.expectedReservations, responseReservations) {
					t.Errorf("Expected %+v; got %+v", tc.expectedReservations, responseReservations)
				}
			}
		})
	}

}

func TestAddReservation(t *testing.T) {
	testCases := []struct {
		testName             string
		booksInDB            []book.Book
		reservationsInDB     []reservation.Reservation
		expectedCode         int
		expectedReservations []reservation.Reservation
		expectedErrorMsg     string
	}{
		// Case 1: GetReservation returns []Reservation
		{
			testName:             "Success",
			booksInDB:            []book.Book{{ID: 1, Title: "Book 1", AvailableCopies: 1}, {ID: 2, Title: "Book 2", AvailableCopies: 2}},
			reservationsInDB:     []reservation.Reservation{{ID: 1, BookId: 1, UserId: 1}},
			expectedCode:         http.StatusCreated,
			expectedReservations: []reservation.Reservation{{ID: 1, BookId: 1, UserId: 1}, {ID: 2, BookId: 2, UserId: 2}},
			expectedErrorMsg:     "",
		},
		// Case 2: GetReservation returns an error
		{
			testName:             "Error",
			booksInDB:            []book.Book{},
			reservationsInDB:     []reservation.Reservation{},
			expectedCode:         http.StatusInternalServerError,
			expectedReservations: nil,
			expectedErrorMsg:     "Could not fetch reservations!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			book_repo := test.NewMockBookRepo(tc.booksInDB)
			res_repo := NewMockReservationRepo(tc.reservationsInDB)
			res_handler := reservation.NewHandler(res_repo, book_repo)

			// HTTP request
			reqBody := `{"book_id": 1}`
			req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Gin context
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(w)
			context.Request = req
			context.Set("userId", 1)

			// Perform the request
			res_handler.AddReservation(context)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}
		})
	}

	// t.Run("Successful Reservation", func(t *testing.T) {

	// 	// HTTP request
	// 	reqBody := `{"BookId": 1}`
	// 	req := httptest.NewRequest(http.MethodPost, "/add-reservation", strings.NewReader(reqBody))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	w := httptest.NewRecorder()

	// 	// Gin context
	// 	gin.SetMode(gin.TestMode)
	// 	context, _ := gin.CreateTestContext(w)
	// 	context.Request = req

	// 	// Perform the request
	// 	handler.AddReservation(context)

	// 	// Assertions
	// 	// assert.Equal(t, http.StatusCreated, w.Code)
	// 	// assert.Contains(t, w.Body.String(), "Reservation added!")
	// })

	// // Add more test cases for other scenarios...

	// // Example: Test case for parsing error
	// t.Run("Parsing Error", func(t *testing.T) {
	// 	// HTTP request with invalid JSON
	// 	req := httptest.NewRequest(http.MethodPost, "/add-reservation", strings.NewReader("invalid JSON"))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	w := httptest.NewRecorder()

	// 	// Gin context
	// 	gin.SetMode(gin.TestMode)
	// 	context, _ := gin.CreateTestContext(w)
	// 	context.Request = req

	// 	// Perform the request
	// 	handler.AddReservation(context)

	// 	// Assertions
	// 	// assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	// assert.Contains(t, w.Body.String(), "Could not parse request data!")
	// })

	// Add more test cases for other scenarios...
}
