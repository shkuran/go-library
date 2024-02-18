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
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reservations!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK, reservations)
}

func addReservation(context *gin.Context) {
	var reservation models.Reservation
	err := context.ShouldBindJSON(&reservation)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		log.Println(err)
		return
	}

	book, err := models.GetBookById(reservation.BookId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch book!"})
		log.Println(err)
		return
	}

	numberOfBookCopies := book.AvailableCopies
	if numberOfBookCopies < 1 {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "The book are not available!"})
		return
	}

	userId := context.GetInt64("userId")
	reservation.UserId = userId

	err = models.SaveReservation(reservation)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not add reservation!"})
		log.Println(err)
		return
	}

	err = models.UpdateAvailableCopies(book.Id, numberOfBookCopies-1)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not update number of book copies!"})
		log.Println(err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Reservation added!"})
}

func returnBook(context *gin.Context) {
	reservationId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not parse reservationId!"})
		log.Println(err)
		return
	}

	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reservation!"})
		log.Println(err)
		return
	}

	userId := context.GetInt64("userId")
	if (reservation.UserId != userId) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not access to copmlete reservation!"})
		return
	}

	if reservation.ReturnDate != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "The reservation is copleted already!"})
		return
	}

	err = models.UpdateReturnDate(reservationId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not copmlete reservation!"})
		log.Println(err)
		return
	}

	book, err := models.GetBookById(reservation.BookId)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch book!"})
		log.Println(err)
		return
	}
	err = models.UpdateAvailableCopies(book.Id, book.AvailableCopies+1)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not update number of book copies!"})
		log.Println(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Reservation copmleted!"})
}
