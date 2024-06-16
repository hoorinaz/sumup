package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// SQL query to create accounts table if it doesn't exist
var createAccountsTableQuery = `
	CREATE TABLE IF NOT EXISTS accounts (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		owner VARCHAR(100) NOT NULL,
		balance FLOAT DEFAULT 0.0
	);`

// SQL query to create transactions table if it doesn't exist
var createTransactionsTableQuery = `
	CREATE TABLE IF NOT EXISTS transactions (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		from_account BIGINT NOT NULL,
		to_account BIGINT NOT NULL,
		amount FLOAT NOT NULL,
		FOREIGN KEY (from_account) REFERENCES accounts(id),
		FOREIGN KEY (to_account) REFERENCES accounts(id)
	);`

func InitDB() (*sql.DB, error) {
	// Database connection setup
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbName

	var db *sql.DB
	var err error

	// Retry mechanism
	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Printf("Failed to connect to database. Attempt %d/%d. Retrying in 5 seconds...", attempts, maxAttempts)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
	}

	// Execute the queries
	if _, err := db.Exec(createAccountsTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create accounts table: %w", err)
	}

	if _, err := db.Exec(createTransactionsTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create transactions table: %w", err)
	}

	return db, nil
}
