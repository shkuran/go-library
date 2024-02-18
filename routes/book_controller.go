package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/models"
)

func getBooks(context *gin.Context) {
	books, err := models.GetBooks()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch books!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK, books)
}

func addBook(context *gin.Context) {
	var user models.Book
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		log.Println(err)
		return
	}
	err = models.SaveBook(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add book!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Book added!"})
}
