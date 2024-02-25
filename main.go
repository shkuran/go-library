package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/book"
	"github.com/shkuran/go-library/db"
	"github.com/shkuran/go-library/reservation"
	"github.com/shkuran/go-library/routes"
	"github.com/shkuran/go-library/user"
)

func main() {

	mysql, err := db.InitDB()
	if err != nil {
		log.Fatalln("Failed to initialize database")
		return
	}

	db.CreateTables(mysql)

	server := gin.Default()

	books_repo := book.NewMySQLRepository(mysql)
	books_controller := book.NewBooksController(books_repo)

	reservation_repo := reservation.NewMySQLRepository(mysql)
	reservation_controller := reservation.NewReservationController(reservation_repo, books_repo)

	user_repo := user.NewMySQLRepository(mysql)
	user_controller := user.NewUserController(user_repo)

	routes.RegisterRoutes(server, reservation_controller, user_controller, books_controller)

	server.Run(":8080")
}
