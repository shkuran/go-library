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
	connStr := "host=host.docker.internal port=5432 user=root password=root dbname=library sslmode=disable"
	driverName := "postgres"

	varDb, err := db.InitDB(driverName, connStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	db.CreateTables(varDb)

	server := gin.Default()

	bookRepo := book.NewRepo(varDb)
	bookHandler := book.NewHandler(bookRepo)

	reservationRepo := reservation.NewRepo(varDb)
	reservationHandler := reservation.NewHandler(reservationRepo, bookRepo)

	userRepo := user.NewRepo(varDb)
	userHandler := user.NewHandler(userRepo)

	routes.RegisterRoutes(server, bookHandler, userHandler, reservationHandler)

	err = server.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
