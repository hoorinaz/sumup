package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
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

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// Execute the queries
	if _, err := db.Exec(createAccountsTableQuery); err != nil {
		log.Fatalf("Failed to create accounts table: %v", err)
	}

	if _, err := db.Exec(createTransactionsTableQuery); err != nil {
		log.Fatalf("Failed to create transactions table: %v", err)
	}
	return db, nil
}
