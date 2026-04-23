package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connection pool to db
var Pool *pgxpool.Pool

func Connect() {
	// Retrieve database URL from .env file
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	// Connect to db
	var err error
	Pool, err = pgxpool.New(context.Background(), db_url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping database
	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not ping database: %v\n", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
