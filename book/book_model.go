package book

type Book struct {
	ID              int64  `json:"id" db:"id"`
	Title           string `json:"title" db:"title" binding:"required"`
	Author          string `json:"author" db:"author" binding:"required"`
	ISBN            string `json:"isbn" db:"isbn" binding:"required"`
	PublicationYear int64  `json:"publication_year" db:"publication_year" binding:"required"`
	AvailableCopies int64  `json:"available_copies" db:"available_copies" binding:"required"`
}