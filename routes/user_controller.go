package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/models"
	"github.com/shkuran/go-library/utils"
)

func createUser(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		log.Println(err)
		return
	}
	err = models.SaveUser(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}

func getUsers(context *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch books!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data!"})
		log.Println(err)
		return
	}

	err = models.ValidateCredentials(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user!"})
		log.Println(err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token!"})
		log.Println(err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successfully!", "token": token})
}