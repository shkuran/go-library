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
	connStr := "host=localhost port=5432 user=root password=root dbname=library sslmode=disable"
	driverName := "postgres"

	varDb, err := db.InitDB(driverName, connStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	db.CreateTables(varDb)

	server := gin.Default()

	book_repo := book.NewRepo(varDb)
	book_handler := book.NewHandler(book_repo)

	reservation_repo := reservation.NewRepo(varDb)
	reservation_handler := reservation.NewHandler(reservation_repo, book_repo)

	user_repo := user.NewRepo(varDb)
	user_handler := user.NewHandler(user_repo)

	routes.RegisterRoutes(server, book_handler, user_handler, reservation_handler)

	err = server.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

// host := "localhost"
// port := 5432
// user := "root"
// password := "root"
// dbname := "library"

// connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
