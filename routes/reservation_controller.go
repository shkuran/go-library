package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/models"
)

func getReservations(context *gin.Context) {
	reservations, err := models.GetReservations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reservations!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK, reservations)
}

func addReservation(context *gin.Context) {
	var reservation models.Reservation
	err := context.ShouldBindJSON(&reservation)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		log.Println(err)
		return
	}

	//TODO: add check book and client id
	// reduce copy of books

	id, err := models.AddReservation(reservation)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add reservation!"})
		log.Println(err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Reservation added!", "id": id})
}

func returnBook(context *gin.Context) {
	reservationId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse reservationId!"})
		log.Println(err)
		return
	}

	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reservation!"})
		log.Println(err)
		return
	}
	if reservation.ReturnDate != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Reservation is copleted already!"})
		return
	}

	err = models.UpdateReturnDate(reservationId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not copmlete reservation!"})
		log.Println(err)
		return
	}

	//TODO: increase number of books

	context.JSON(http.StatusOK, gin.H{"message": "Reservation copmleted!"})
}
