package test

import (
	"errors"
	"time"

	"github.com/shkuran/go-library/reservation"
)

type MockReservationRepo struct {
	reservation []reservation.Reservation
}

func NewMockReservationRepo(res []reservation.Reservation) *MockReservationRepo {
	return &MockReservationRepo{reservation: res}
}

func (r *MockReservationRepo) GetAll() ([]reservation.Reservation, error) {
	if len(r.reservation) == 0 {
		return nil, errors.New("simulated error fetching reservations")
	}
	return r.reservation, nil
}

func (r *MockReservationRepo) GetById(id int64) (reservation.Reservation, error) {
	for _, res := range r.reservation {
		if res.ID == id {
			return res, nil
		}
	}
	return reservation.Reservation{}, errors.New("simulated error fetching reservation by id")
}

func (r *MockReservationRepo) Save(res reservation.Reservation) error {
	r.reservation = append(r.reservation, res)
	return nil
}

func (r *MockReservationRepo) UpdateReturnDate(id int64) error {
	returnDate := time.Now()
	r.reservation[id].ReturnDate = &returnDate
	return nil
}