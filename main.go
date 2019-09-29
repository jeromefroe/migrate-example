package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Migrate database to desired schema if need be.
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("could not construct new database migration: %v", err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("could not migrate database: %v", err)
		}
		log.Print("no change needed for database")
	}

	// Open a connection to the database.
	db, err := sql.Open("postgres", "user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer db.Close()

	// Insert a record into the database.
	insertStmt, err := db.Prepare("INSERT INTO users(username,password,email) VALUES($1, $2, $3) returning user_id;")
	if err != nil {
		log.Fatalf("could not prepare insert statement: %v", err)
	}
	defer insertStmt.Close()

	insertRows, err := insertStmt.Query("jerome", "password", "jerome@example.com")
	if err != nil {
		log.Fatalf("could not execute insert statement: %v", err)
	}
	defer insertRows.Close()

	if !insertRows.Next() {
		log.Fatal("could not get id of inserted record")
	}

	var id int
	if err := insertRows.Scan(&id); err != nil {
		log.Fatalf("could not get id of inserted record: %v", err)
	}
	fmt.Println("inserted record id = ", id)

	// Select a specific record from the database.
	whereStmt, err := db.Prepare("SELECT user_id, username, email FROM users WHERE user_id=$1;")
	if err != nil {
		log.Fatalf("could not prepare where statement: %v", err)
	}

	whereRows, err := whereStmt.Query(id)
	if err != nil {
		log.Fatalf("could not execute where statement: %v", err)
	}
	defer whereRows.Close()

	if !whereRows.Next() {
		log.Fatal("no records found for inserted id")
	}

	var (
		user_id         int
		username, email string
	)
	if err := whereRows.Scan(&user_id, &username, &email); err != nil {
		log.Fatalf("could not scan fields: %v", err)
	}
	fmt.Printf("user_id: %v, username: %v, email: %v\n", user_id, username, email)

	// Select multiple rows from the database.
	selectStmt, err := db.Prepare("SELECT user_id, username, email FROM users;")
	if err != nil {
		log.Fatalf("could not prepare select statement: %v", err)
	}

	selectRows, err := selectStmt.Query()
	if err != nil {
		log.Fatalf("could not execute select statement: %v", err)
	}
	defer selectRows.Close()

	for selectRows.Next() {
		var (
			user_id         int
			username, email string
		)
		if err := selectRows.Scan(&user_id, &username, &email); err != nil {
			log.Fatalf("could not scan fields: %v", err)
		}
		fmt.Printf("user_id: %v, username: %v, email: %v\n", user_id, username, email)
	}
}
