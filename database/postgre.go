package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
	dsn := "host=localhost port=5432 user=postgres password=raisyaadmin dbname=prestasi_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL connected")
	return db
}
