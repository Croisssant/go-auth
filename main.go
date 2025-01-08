package main

import (
	"net/http"

	"croissant.com/go/auth/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	basicAuth := router.Group("/basic-auth")
	{
		basicAuth.GET("/ping",
			auth.BasicAuthHandlerFunc(),
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
	}

	bearerAuth := router.Group("/bearer")
	{
		bearerAuth.POST("/login", auth.BasicAuthHandlerFunc(), auth.BearerTokenGen)
		bearerAuth.GET("/ping", auth.BearerTokenCheck)
	}

	jwtAuth := router.Group("/jwt")
	{
		jwtAuth.POST("/login", auth.BasicAuthHandlerFunc(), auth.JwtTokenGen)
		jwtAuth.GET("/ping", auth.JwtTokenCheck)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}
