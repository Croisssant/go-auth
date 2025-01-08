package auth

import "github.com/gin-gonic/gin"

func BasicAuthHandlerFunc() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	})
}
