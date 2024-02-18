package models

import (
	"time"

	"github.com/shkuran/go-library/db"
)

type Reservation struct {
	Id           int64      `json:"id" db:"id"`
	BookId       int64      `json:"book_id" db:"book_id"`
	UserId       int64      `json:"user_id" db:"user_id"`
	CheckoutDate time.Time  `json:"checkout_date" db:"checkout_date"`
	ReturnDate   *time.Time `json:"return_date" db:"return_date"`
}

func GetReservations() ([]Reservation, error) {
	rows, err := db.DB.Query("SELECT * FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var res Reservation
		err := rows.Scan(&res.Id, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}

func GetReservationById(id int64) (Reservation, error) {
	var res Reservation

	row := db.DB.QueryRow("SELECT * FROM reservations WHERE id = ?", id)
	err := row.Scan(&res.Id, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
	if err != nil {
		return res, err
	}

	return res, nil
}

func SaveReservation(res Reservation) error {
	query := `
	INSERT INTO reservations (book_id, user_id, checkout_date) 
	VALUES (?, ?, ?)
	`
	reservationDate := time.Now()

	_, err := db.DB.Exec(query, res.BookId, res.UserId, reservationDate)
	if err != nil {
		return err
	}

	return nil
}

func UpdateReturnDate(id int64) error {
	query := `
	UPDATE reservations
	SET return_date = ?
	WHERE id = ?
	`
	returnDate := time.Now()

	_, err := db.DB.Exec(query, returnDate, id)

	return err
}
