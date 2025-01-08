package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Bearer Token Based Authentication
var bearerTokens []string

func BearerTokenGen(ctx *gin.Context) {
	token, _ := randomHex(20)
	bearerTokens = append(bearerTokens, token)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func BearerTokenCheck(ctx *gin.Context) {

	bearerToken := ctx.Request.Header.Get("Authorization")
	reqToken := strings.Split(bearerToken, " ")[1]
	for _, token := range bearerTokens {
		if token == reqToken {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
			return
		}
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"message": "Unauthorized: Access Denied",
	})

}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
