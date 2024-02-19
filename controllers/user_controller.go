package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/models"
	"github.com/shkuran/go-library/utils"
)

func CreateUser(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = models.SaveUser(user)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not create user!", err)
		return
	}
	utils.HandleStatusCreated(context, "User created!")
}

func GetUsers(context *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch users!", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func Login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}

	err = models.ValidateCredentials(&user)
	if err != nil {
		utils.HandleStatusUnauthorized(context, "Could not authenticate user!", err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not generate token!", err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successfully!", "token": token})
}
