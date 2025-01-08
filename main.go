package main

import (
	"net/http"

	"croissant.com/go/auth/auth"
	"croissant.com/go/auth/models"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

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

	publicRoutes := router.Group("/public")
	{
		publicRoutes.POST("/login", auth.Login)
		publicRoutes.POST("/register", auth.Register)
	}

	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(auth.AuthenticationMiddleware())
	{
		protectedRoutes.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		protectedRoutes.GET("/albums", models.GetAlbums)
		protectedRoutes.GET("/albums/:id", models.GetAlbumById)
		protectedRoutes.POST("/albums", models.PostAlbums)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}
