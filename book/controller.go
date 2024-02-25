package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/utils"
)

type BooksController struct {
	repo Repository
}

func NewBooksController(repo Repository) BooksController {
	return BooksController{repo: repo}
}

func (svc BooksController) GetBooks(context *gin.Context) {
	books, err := svc.repo.GetBooks()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch books!", err)
		return
	}
	context.JSON(http.StatusOK, books)
}

func (svc BooksController) AddBook(context *gin.Context) {
	var b Book
	err := context.ShouldBindJSON(&b)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = svc.repo.SaveBook(b)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not add book!", err)
		return
	}
	utils.HandleStatusCreated(context, "Book added!")
}
