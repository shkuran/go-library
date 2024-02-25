package reservation

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	Repository{db: db}
}

func (r *Repository) GetReservations() ([]Reservation, error) {
	return r.db.GetReservations()
}

func (r *Repository) GetReservationById(id int64) (Reservation, error) {
	return r.db.GetReservationById(id)
}

func (r *Repository) SaveReservation(res Reservation) error {
	return r.db.SaveReservation(res)
}

func (r *Repository) UpdateReturnDate(id int64) error {
	return r.db.UpdateReturnDate(id)
}
