package user

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

func (h Handler) CreateUser(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = h.repo.save(user)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not create user!", err)
		return
	}
	utils.HandleStatusCreated(context, "User created!")
}

func (h Handler) GetUsers(context *gin.Context) {
	users, err := h.repo.getAll()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch users!", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func (h Handler) Login(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}

	err = h.repo.validateCredentials(&user)
	if err != nil {
		utils.HandleStatusUnauthorized(context, "Could not authenticate user!", err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not generate token!", err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successfully!", "token": token})
}
