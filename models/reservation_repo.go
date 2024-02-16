package models

import (
	"time"

	"github.com/shkuran/go-library/db"
)

type Reservation struct {
	Id           int64
	BookId       int64
	ClientId     int64
	CheckoutDate time.Time
	ReturnDate   *time.Time
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
		err := rows.Scan(&res.Id, &res.BookId, &res.ClientId, &res.CheckoutDate, &res.ReturnDate)
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
	err := row.Scan(&res.Id, &res.BookId, &res.ClientId, &res.CheckoutDate, &res.ReturnDate)
		if err != nil {
			return res, err
		}

	return res, nil
}

func AddReservation(res Reservation) (int64, error) {
	query := `
	INSERT INTO reservations (book_id, client_id, checkout_date) 
	VALUES (?, ?, ?)
	`
	reservationDate := time.Now()

	result, err := db.DB.Exec(query, res.BookId, res.ClientId, reservationDate)
	if err != nil {
		return 0, err
	}

	reservationId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return reservationId, nil
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
