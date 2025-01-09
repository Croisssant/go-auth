package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Try this next https://github.com/gin-gonic/gin/issues/932
// Find ways to use goroutines and channel

func DBInit() (*sql.DB, error) {
	// Get database URL and auth token from environment variables
	dbUrl := os.Getenv("TURSO_DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("TURSO_URL environment variable not set")
	}

	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if authToken != "" {
		dbUrl += "?authToken=" + authToken
	}

	// Open database connection
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error opening cloud db: %w", err)
	}

	return db, nil
}

func DBMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("DB", db)
		ctx.Next()
	}
}

func DBRetrieve(ctx *gin.Context) {
	db := ctx.MustGet("DB").(*sql.DB)
	// Query the data
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		fmt.Printf("error querying data: %v", err)
	}
	defer rows.Close()

	// Print results
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			fmt.Printf("error scanning row: %v", err)
		}
		fmt.Printf("Row: id=%d, name=%s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error iterating rows: %v", err)
	}
}

// func run() (err error) {

// 	defer db.Close()

// 	// Configure connection pool
// 	db.SetConnMaxIdleTime(9 * time.Second)

// 	// Create test table
// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
// 	if err != nil {
// 		return fmt.Errorf("error creating table: %w", err)
// 	}

// 	// Check if test data already exists
// 	var exists bool
// 	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM test WHERE id = 1)").Scan(&exists)
// 	if err != nil {
// 		return fmt.Errorf("error checking existing data: %w", err)
// 	}

// 	// Insert test data only if it doesn't exist
// 	if !exists {
// 		_, err = db.Exec("INSERT INTO test (id, name) VALUES (?, ?)", 1, "remote test")
// 		if err != nil {
// 			return fmt.Errorf("error inserting data: %w", err)
// 		}
// 		fmt.Println("Inserted test data")
// 	} else {
// 		fmt.Println("Test data already exists, skipping insert")
// 	}

// 	// Query the data
// 	rows, err := db.Query("SELECT * FROM test")
// 	if err != nil {
// 		return fmt.Errorf("error querying data: %w", err)
// 	}
// 	defer rows.Close()

// 	// Print results
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		if err := rows.Scan(&id, &name); err != nil {
// 			return fmt.Errorf("error scanning row: %w", err)
// 		}
// 		fmt.Printf("Row: id=%d, name=%s\n", id, name)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return fmt.Errorf("error iterating rows: %w", err)
// 	}

// 	fmt.Printf("Successfully connected and executed queries on %s\n", dbUrl)
// 	return nil
// }
