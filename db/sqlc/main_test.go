package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Load environment variables from the .env file
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatalf("Error loading .env file: %v", env_err)
	}

	// Fetch driver and source from environment variables
	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")

	// Check if the environment variables were fetched successfully
	if dbDriver == "" {
		log.Fatal("DB_DRIVER not set in environment variables")
	}
	if dbSource == "" {
		log.Fatal("DB_SOURCE not set in environment variables")
	}

	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to Database:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
