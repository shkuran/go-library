package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/utils"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) Handler {
	return Handler{repo: repo}
}

func (h Handler) GetBooks(context *gin.Context) {
	books, err := h.repo.getAll()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch books!", err)
		return
	}
	context.JSON(http.StatusOK, books)
}

// func (h Handler) AddBook(context *gin.Context) {
// 	var b Book
// 	err := context.ShouldBindJSON(&b)
// 	if err != nil {
// 		utils.HandleBadRequest(context, "Could not parse request data!", err)
// 		return
// 	}
// 	err = h.repo.save(b)
// 	if err != nil {
// 		utils.HandleInternalServerError(context, "Could not add book!", err)
// 		return
// 	}
// 	utils.HandleStatusCreated(context, "Book added!")
// }
