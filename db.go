package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

func AddBook(book Book) error {
	_, err := db.Exec(
		`INSERT INTO books (title, author, rating, year, description, type)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		book.Title,
		book.Author,
		book.Rating,
		book.Year,
		book.Description,
		book.Type,
	)

	return err
}

func GetBooks() ([]Book, error) {

	rows, err := db.Query(
		`SELECT title, author, rating, year, description, type FROM books`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book

		err := rows.Scan(
			&b.Title,
			&b.Author,
			&b.Rating,
			&b.Year,
			&b.Description,
			&b.Type,
		)

		if err != nil {
			continue
		}

		books = append(books, b)
	}

	return books, nil
}
