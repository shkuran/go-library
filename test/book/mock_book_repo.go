package test

import (
	"errors"

	"github.com/shkuran/go-library/book"
)

type MockBookRepo struct {
	books []book.Book
}

func NewMockBookRepo(books []book.Book) *MockBookRepo {
	return &MockBookRepo{books: books}
}

func (r *MockBookRepo) GetAll() ([]book.Book, error) {
	if len(r.books) == 0 {
		return nil, errors.New("simulated error fetching books")
	}
	return r.books, nil
}

func (r *MockBookRepo) GetById(id int64) (book.Book, error) {
	for _, bk := range r.books {
		if bk.ID == id {
			return bk, nil
		}
	}
	return book.Book{}, errors.New("simulated error fetching book by id")
}

func (r *MockBookRepo) UpdateAvailableCopies(id, copies int64) error {
	r.books[id].AvailableCopies = copies
	return nil
}

func (r *MockBookRepo) Save(book book.Book) error {
	r.books = append(r.books, book)
	return nil
}
