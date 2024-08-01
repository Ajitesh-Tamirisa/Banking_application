package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Load .env file in non-CI environments
func init() {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		}
	}
}

var (
	dbDriver = os.Getenv("DB_DRIVER")
	dbSource = os.Getenv("DB_SOURCE")
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	if dbDriver == "" || dbSource == "" {
		log.Fatal("Database environment variables not set")
	}

	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to Database:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
