package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/auth"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/reservations", getReservations)
	server.POST("/reservations", auth.Authenticate, addReservation)
	server.POST("/reservations/:id", auth.Authenticate, returnBook)

	server.GET("/books", getBooks)
	server.POST("/books", addBook)

	server.POST("/signup", createUser)
	server.POST("/login", login)
	server.GET("/users", getUsers)
}
