package conifg

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var DB *sql.DB
var err error

// DBConfig holds the database connection details
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// LoadDBConfig loads the database configuration (could be extended to load from env or config files)
func LoadDBConfig() DBConfig {
	// For demo purposes, hardcoding here, but ideally use environment variables or config files
	return DBConfig{
		Host:     "localhost",     // MySQL host
		Port:     "3306",          // MySQL port
		Username: "phpmyadmin",    // MySQL username
		Password: "your_password", // MySQL password
		Database: "movie_rds",     // MySQL database name
	}
}

// DatabaseConnection initializes the DB connection and applies connection pooling settings
func DatabaseConnection() {
	// Load configuration
	config := LoadDBConfig()

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// Open a connection to MySQL using sql.DB
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Test the database connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(100)          // Set the maximum number of open connections
	DB.SetMaxIdleConns(10)           // Set the maximum number of idle connections
	DB.SetConnMaxLifetime(time.Hour) // Set the maximum lifetime of a connection

	// Optional: Log the success message
	fmt.Println("Database connected successfully...")
}

// CloseDBConnection closes the database connection
func CloseDBConnection() {
	if err := DB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		fmt.Println("Database connection closed.")
	}
}
