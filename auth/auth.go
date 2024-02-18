package auth

import (
	"net/http"

	"github.com/shkuran/go-library/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Empty token!"})
		return
	}

	userId, err := utils.VerifyTokenAndReturnUserId(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token!"})
		return
	}

	context.Set("userId", userId)
	context.Next()
}