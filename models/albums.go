package models

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, albums)
}

func PostAlbums(ctx *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum
	if err := ctx.BindJSON(&newAlbum); err != nil {
		fmt.Printf("Error: %v \n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid format"})
		ctx.Abort()
		return
	}

	// Add the new album to the slice
	albums = append(albums, newAlbum)
	ctx.JSON(http.StatusCreated, newAlbum)
}

func GetAlbumById(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, alb := range albums {
		if alb.ID == id {
			ctx.JSON(http.StatusOK, alb)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
