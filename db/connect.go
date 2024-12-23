package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	ID        uint32
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}

func ConnectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=fetcher password=dev dbname=mentoref sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to database")

	return db
}
