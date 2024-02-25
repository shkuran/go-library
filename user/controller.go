package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shkuran/go-library/utils"
)

type UserController interface {
	CreateUser(context *gin.Context)
	GetUsers(context *gin.Context)
	Login(context *gin.Context)
}

type UserControllerImpl struct {
	repo Repository
}

func NewUserController(repo Repository) *UserControllerImpl {
	return &UserControllerImpl{repo: repo}
}

func (ctr UserControllerImpl) CreateUser(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}
	err = ctr.repo.SaveUser(user)
	if err != nil {
		utils.HandleInternalServerError(context, "Could not create user!", err)
		return
	}
	utils.HandleStatusCreated(context, "User created!")
}

func (ctr UserControllerImpl) GetUsers(context *gin.Context) {
	users, err := ctr.repo.GetUsers()
	if err != nil {
		utils.HandleInternalServerError(context, "Could not fetch users!", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func (ctr UserControllerImpl) Login(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.HandleBadRequest(context, "Could not parse request data!", err)
		return
	}

	err = ctr.repo.ValidateCredentials(&user)
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