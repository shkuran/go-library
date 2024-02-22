package book

func GetBookById(bookId int64) (Book, error) {
	b, err := getBookById(bookId)
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

func UpdateNumberOfBooks(bookId, copies int64) error {
	err := updateAvailableCopies(bookId, copies)
	return err
}