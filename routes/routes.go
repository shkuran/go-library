package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/auth"
	"github.com/shkuran/go-library/book"
	"github.com/shkuran/go-library/reservation"
	"github.com/shkuran/go-library/user"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/reservations", reservation.GetReservations)
	server.POST("/reservations", auth.Authenticate, reservation.AddReservation)
	server.POST("/reservations/:id", auth.Authenticate, reservation.CopleteReservation)

	server.GET("/books", func(context *gin.Context) {
        book.GetBooks(context, book.GetBooksFromDB)
    })
	server.POST("/books", book.AddBook)

	server.POST("/signup", user.CreateUser)
	server.POST("/login", user.Login)
	server.GET("/users", user.GetUsers)
}
