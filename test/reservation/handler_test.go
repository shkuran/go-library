package test

import (
	"encoding/json"
	"fmt"
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

func setupTestEnv(booksInDB []book.Book, reservationsInDB []reservation.Reservation) TestEnv {
	bookRepo := test.NewMockBookRepo(booksInDB)
	resRepo := NewMockReservationRepo(reservationsInDB)
	resHandler := reservation.NewHandler(resRepo, bookRepo)

	return TestEnv{
		BookRepo:           bookRepo,
		ReservationRepo:    resRepo,
		ReservationHandler: resHandler,
	}
}

type TestEnv struct {
	BookRepo           *test.MockBookRepo
	ReservationRepo    *MockReservationRepo
	ReservationHandler reservation.Handler
}

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
			// Set up the test environment
			env := setupTestEnv(tc.booksInDB, tc.reservationsInDB)

			router := gin.Default()

			router.GET("/reservations", env.ReservationHandler.GetReservations)

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
				// Check if the response matches the expected reservations
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
		testName         string
		booksInDB        []book.Book
		reservationsInDB []reservation.Reservation
		requestBody      string
		expectedCode     int
		expectedErrorMsg string
	}{
		// Case 1: AddReservation adds new reservation and update AvailableCopies
		{
			testName:         "Successfull added reservation",
			booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []reservation.Reservation{},
			requestBody:      `{"book_id": 1}`,
			expectedCode:     http.StatusCreated,
			expectedErrorMsg: "",
		},
		// Case 2: AddReservation returns a bad request
		{
			testName:         "Bad request",
			booksInDB:        []book.Book{},
			reservationsInDB: []reservation.Reservation{},
			requestBody:      `{"book_id": 1a}`,
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "Could not parse request data!",
		},
		// Case 3: AddReservation could not fetch book! Returns InternalServerError
		{
			testName:         "No books",
			booksInDB:        []book.Book{},
			reservationsInDB: []reservation.Reservation{},
			requestBody:      `{"book_id": 18}`,
			expectedCode:     http.StatusInternalServerError,
			expectedErrorMsg: "Could not fetch book!",
		},
		// Case 4: AddReservation returns a bad request. The book is not available!
		{
			testName:         "AvailableCopies is 0",
			booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 0}},
			reservationsInDB: []reservation.Reservation{},
			requestBody:      `{"book_id": 1}`,
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "The book is not available!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Set up the test environment
			env := setupTestEnv(tc.booksInDB, tc.reservationsInDB)

			// HTTP request
			req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Gin context
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(w)
			context.Request = req
			context.Set("userId", int64(1))
			context.Next()

			userId := context.GetInt64("userId")
			fmt.Println(userId)

			// Perform the request
			env.ReservationHandler.AddReservation(context)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}

			if tc.expectedErrorMsg == "" {
				// Check if AvailableCopies of book with id:1 was updated(was: 1, should be: 0)
				reservedBook, err := env.BookRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch book! error: %d", err)
				}
				if reservedBook.AvailableCopies != 0 {
					t.Errorf("Expected AvailableCopies %d; got %d", 0, reservedBook.AvailableCopies)
				}
				// Check if reservation was added
				expRes := reservation.Reservation{ID: 1, BookId: 1, UserId: 1}
				gotedRes, err := env.ReservationRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch book! error: %d", err)
				}
				if !reflect.DeepEqual(gotedRes, expRes) {
					t.Errorf("Expected new rservation id:%d, book_id:%d, user_id:%d; got id:%d, book_id:%d, user_id:%d",
						expRes.ID, expRes.BookId, expRes.UserId, gotedRes.ID, gotedRes.BookId, gotedRes.UserId)
				}
			} else {
				// Check if the response contains the expected error message
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}

			}
		})
	}

}

func TestCompleteReservation(t *testing.T) {
	testCases := []struct {
		testName         string
		booksInDB        []book.Book
		reservationsInDB []reservation.Reservation
		reservationId    int
		expectedCode     int
		expectedErrorMsg string
	}{
		// Case 1: CopleteReservation add return date for reservation and update AvailableCopies for book
		{
			testName:         "Successfull completed reservation",
			booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 1}},
			reservationsInDB: []reservation.Reservation{{ID: 1, BookId: 1, UserId: 1, ReturnDate: nil}},
			expectedCode:     http.StatusOK,
			expectedErrorMsg: "",
		},
		// Case 2: CopleteReservation returns a bad request
		// {
		// 	testName:         "Bad request",
		// 	booksInDB:        []book.Book{},
		// 	reservationsInDB: []reservation.Reservation{},
		// 	expectedCode:     http.StatusBadRequest,
		// 	expectedErrorMsg: "Could not parse request data!",
		// },
		// // Case 3: CopleteReservation could not fetch book! Returns InternalServerError
		// {
		// 	testName:         "No books",
		// 	booksInDB:        []book.Book{},
		// 	reservationsInDB: []reservation.Reservation{},
		// 	expectedCode:     http.StatusInternalServerError,
		// 	expectedErrorMsg: "Could not fetch book!",
		// },
		// // Case 4: CopleteReservation returns a bad request. The book is not available!
		// {
		// 	testName:         "AvailableCopies is 0",
		// 	booksInDB:        []book.Book{{ID: 1, Title: "Book_1", AvailableCopies: 0}},
		// 	reservationsInDB: []reservation.Reservation{},
		// 	expectedCode:     http.StatusBadRequest,
		// 	expectedErrorMsg: "The book is not available!",
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// Set up the test environment
			env := setupTestEnv(tc.booksInDB, tc.reservationsInDB)

			// HTTP request
			req := httptest.NewRequest(http.MethodPost, "/reservations", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Gin context
			gin.SetMode(gin.TestMode)
			context, _ := gin.CreateTestContext(w)
			context.Request = req
			context.Set("userId", int64(1))
			context.AddParam("id", "1")

			// Perform the request
			env.ReservationHandler.CopleteReservation(context)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d; got %d", tc.expectedCode, w.Code)
			}

			if tc.expectedErrorMsg == "" {
				// Check if AvailableCopies of book with id:1 was updated(was: 1, should be: 2)
				reservedBook, err := env.BookRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch book! error: %d", err)
				}
				if reservedBook.AvailableCopies != 2 {
					t.Errorf("Expected AvailableCopies %d; got %d", 2, reservedBook.AvailableCopies)
				}
				// Check if reservation was added
				expRes := reservation.Reservation{ID: 1, BookId: 1, UserId: 1}
				gotedRes, err := env.ReservationRepo.GetById(1)
				if err != nil {
					t.Errorf("Could not fetch reservation! error: %d", err)
				}
				if !reflect.DeepEqual(gotedRes, expRes) {
					t.Errorf("Expected new rservation: %v; got: %v", expRes, gotedRes)
					// t.Errorf("Expected new rservation id:%d, book_id:%d, user_id:%d; got id:%d, book_id:%d, user_id:%d",
					// 	expRes.ID, expRes.BookId, expRes.UserId, gotedRes.ID, gotedRes.BookId, gotedRes.UserId)
				}
			} else {
				// Check if the response contains the expected error message
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}
				if response["message"] != tc.expectedErrorMsg {
					t.Errorf("Expected error message '%s'; got '%s'", tc.expectedErrorMsg, response["message"])
				}

			}
		})
	}

}
