import (
	"database/sql"
	"time"
)

type DB interface {
	GetReservations() ([]Reservation, error)
	GetReservationById(id int64) (Reservation, error)
	SaveReservation(res Reservation) error
	UpdateReturnDate(id int64) error
}

type InMemory struct {
	m map[int64]Reservation
}

func NewInMemory() *InMemory {
	MySQL{m: make(map[int64]Reservation)}
}

func (in_memory *InMemory) GetReservations() ([]Reservation, error) {
	v := make([]string, 0, len(m))

	for _, value := range m {
		v = append(v, value)
	}

	return v
}

// implemente everything else

type MySQL struct {
	db *sql.DB
}

func NewMySQL(db *sql.DB) *MySQL {
	MySQL{db: db}
}

func (mysql *MySQL) GetReservations() ([]Reservation, error) {
	rows, err := mysql.db.Query("SELECT * FROM reservations")
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

func (db *MySQL) GetReservationById(id int64) (Reservation, error) {
	var res Reservation

	row := mysql.db.QueryRow("SELECT * FROM reservations WHERE id = ?", id)
	err := row.Scan(&res.ID, &res.BookId, &res.UserId, &res.CheckoutDate, &res.ReturnDate)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (mysql *MySQL) SaveReservation(res Reservation) error {
	query := `
	INSERT INTO reservations (book_id, user_id, checkout_date) 
	VALUES (?, ?, ?)
	`
	reservationDate := time.Now()

	_, err := mysql.db.Exec(query, res.BookId, res.UserId, reservationDate)
	if err != nil {
		return err
	}

	return nil
}

func (mysql *MySQL) UpdateReturnDate(id int64) error {
	query := `
	UPDATE reservations
	SET return_date = ?
	WHERE id = ?
	`
	returnDate := time.Now()

	_, err := mysql.db.Exec(query, returnDate, id)

	return err
}