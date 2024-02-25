package reservation

import (
	"database/sql"
	"time"
)

type Repository interface {
	GetReservations() ([]Reservation, error)
	GetReservationById(id int64) (Reservation, error)
	SaveReservation(res Reservation) error
	UpdateReturnDate(id int64) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetReservations() ([]Reservation, error) {
	rows, err := r.db.Query("SELECT * FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var res Reservation
		err := rows.Scan(&res.ID, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}

func (r *MySQLRepository) GetReservationById(id int64) (Reservation, error) {
	var res Reservation

	row := r.db.QueryRow("SELECT * FROM reservations WHERE id = ?", id)
	err := row.Scan(&res.ID, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *MySQLRepository) SaveReservation(res Reservation) error {
	query := `
	INSERT INTO reservations (book_id, user_id, checkout_date) 
	VALUES (?, ?, ?)
	`
	reservationDate := time.Now()

	_, err := r.db.Exec(query, res.BookId, res.UserId, reservationDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepository) UpdateReturnDate(id int64) error {
	query := `
	UPDATE reservations
	SET return_date = ?
	WHERE id = ?
	`
	returnDate := time.Now()

	_, err := r.db.Exec(query, returnDate, id)

	return err
}
