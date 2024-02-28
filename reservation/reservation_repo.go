package reservation

import (
	"database/sql"
	"time"

	"github.com/shkuran/go-library/db"
)

type Repository interface {
	getAll() ([]Reservation, error)
	getById(id int64) (Reservation, error)
	save(res Reservation) error
	updateReturnDate(id int64) error
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) getAll() ([]Reservation, error) {
	query := "SELECT * FROM reservations"
	rows, err := db.DB.Query(query)
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

func (r *Repo) getById(id int64) (Reservation, error) {
	var res Reservation
	query := `
	SELECT * FROM reservations 
	WHERE id = $1
	`
	row := db.DB.QueryRow(query, id)
	err := row.Scan(&res.ID, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *Repo) save(res Reservation) error {
	query := `
	INSERT INTO reservations (book_id, user_id, checkout_date) 
	VALUES ($1, $2, $3)
	`
	reservationDate := time.Now()

	_, err := db.DB.Exec(query, res.BookId, res.UserId, reservationDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) updateReturnDate(id int64) error {
	query := `
	UPDATE reservations
	SET return_date = $1
	WHERE id = $2
	`
	returnDate := time.Now()

	_, err := db.DB.Exec(query, returnDate, id)

	return err
}
