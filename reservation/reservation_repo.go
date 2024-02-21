package reservation

import (
	"time"

	"github.com/shkuran/go-library/db"
)

func getReservations() ([]Reservation, error) {
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

func getReservationById(id int64) (Reservation, error) {
	var res Reservation

	row := db.DB.QueryRow("SELECT * FROM reservations WHERE id = ?", id)
	err := row.Scan(&res.Id, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
	if err != nil {
		return res, err
	}

	return res, nil
}

func saveReservation(res Reservation) error {
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

func updateReturnDate(id int64) error {
	query := `
	UPDATE reservations
	SET return_date = ?
	WHERE id = ?
	`
	returnDate := time.Now()

	_, err := db.DB.Exec(query, returnDate, id)

	return err
}
