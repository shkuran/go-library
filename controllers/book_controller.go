package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/models"
	"github.com/shkuran/go-library/utils"
)

func GetBooks(context *gin.Context) {
	books, err := models.GetBooks()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch books!", err)
		return
	}
	context.JSON(http.StatusOK, books)
}

func AddBook(context *gin.Context) {
	var user models.Book
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = models.SaveBook(user)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not add book!", err)
		return
	}
	utils.HandleStatusCreated(context, "Book added!")
}
