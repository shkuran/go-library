package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/reservations", getReservations)
	server.POST("/reservations", addReservation)
	server.POST("/reservations/:id", returnBook)
}
