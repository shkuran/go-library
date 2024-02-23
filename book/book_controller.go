package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/utils"
)

func GetBooks(context *gin.Context, getBooks func() ([]Book, error)) {
	books, err := getBooks()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch books!", err)
		return
	}
	context.JSON(http.StatusOK, books)
}

func AddBook(context *gin.Context) {
	var b Book
	err := context.ShouldBindJSON(&b)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = saveBook(b)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not add book!", err)
		return
	}
	utils.HandleStatusCreated(context, "Book added!")
}
