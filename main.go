package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"croissant.com/go/auth/auth"
	"croissant.com/go/auth/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
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

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running example: %v\n", err)

	}

	router.Run() // listen and serve on 0.0.0.0:8080
}

func run() (err error) {
	// Get database URL and auth token from environment variables
	dbUrl := os.Getenv("TURSO_DATABASE_URL")
	if dbUrl == "" {
		return fmt.Errorf("TURSO_URL environment variable not set")
	}

	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if authToken != "" {
		dbUrl += "?authToken=" + authToken
	}

	// Open database connection
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return fmt.Errorf("error opening cloud db: %w", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetConnMaxIdleTime(9 * time.Second)

	// Create test table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	// Check if test data already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM test WHERE id = 1)").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existing data: %w", err)
	}

	// Insert test data only if it doesn't exist
	if !exists {
		_, err = db.Exec("INSERT INTO test (id, name) VALUES (?, ?)", 1, "remote test")
		if err != nil {
			return fmt.Errorf("error inserting data: %w", err)
		}
		fmt.Println("Inserted test data")
	} else {
		fmt.Println("Test data already exists, skipping insert")
	}

	// Query the data
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		return fmt.Errorf("error querying data: %w", err)
	}
	defer rows.Close()

	// Print results
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		fmt.Printf("Row: id=%d, name=%s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	fmt.Printf("Successfully connected and executed queries on %s\n", dbUrl)
	return nil
}
