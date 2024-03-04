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
	for i := range r.books {
		if r.books[i].ID == id {
			r.books[i].AvailableCopies = copies
			return nil
		}
	}
	return errors.New("simulated error updating AvailableCopies")
}

func (r *MockBookRepo) Save(book book.Book) error {
	book.ID = int64(len(r.books)) + 1
	r.books = append(r.books, book)
	return nil
}
