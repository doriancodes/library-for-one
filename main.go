package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type BookType string

const (
	Read   BookType = "READ"
	ToRead BookType = "TO_READ"
	Rec    BookType = "RECOMMENDATION"
)

type Book struct {
	Title       string
	Author      string
	Rating      int
	Year        int
	Description string
	Type        BookType
}

type PageData struct {
	Read            []Book
	ToRead          []Book
	Recommendations []Book
}

var db *sql.DB
var books []Book

func initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title TEXT,
		author TEXT,
		rating INT,
		year INT,
		description TEXT,
		type TEXT
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	connStr := "host=db port=5432 user=postgres password=postgres dbname=library sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	initDB()

	tmpl, err := template.ParseFiles("static/page.html")
	if err != nil {
		log.Fatal("Template parse error:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// data := PageData{
		// 	Read: []Book{
		// 		{"The Hobbit", "J.R.R. Tolkien", 5, 1937},
		// 		{"1984", "George Orwell", 4, 1949},
		// 	},
		// 	UpNext: []Book{
		// 		{"Dune", "Frank Herbert", 0, 1965},
		// 		{"The Pragmatic Programmer", "Andrew Hunt", 0, 1999},
		// 	},
		// 	Recommendations: []Book{
		// 		{"Brave New World", "Aldous Huxley", 0, 1932},
		// 		{"The Name of the Wind", "Patrick Rothfuss", 0, 2007},
		// 	},
		// }

		// err := tmpl.Execute(w, data)
		// if err != nil {
		// 	log.Println("Execute error:", err)
		// }
		//
		var readBooks []Book
		var toReadBooks []Book
		var recBooks []Book

		rows, err := db.Query("SELECT title, author, rating, year, description, type FROM books")
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

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
				log.Println(err)
				continue
			}

			if b.Type == Read {
				readBooks = append(readBooks, b)
			}
			if b.Type == ToRead {
				toReadBooks = append(toReadBooks, b)
			}
			if b.Type == Rec {
				recBooks = append(recBooks, b)
			}
		}

		data := PageData{
			Read:            readBooks,
			ToRead:          toReadBooks,
			Recommendations: recBooks,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("Template error:", err)
		}
	})

	// FORM SUBMIT
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		rating, _ := strconv.Atoi(r.FormValue("rating"))
		year, _ := strconv.Atoi(r.FormValue("year"))

		book := Book{
			Title:       r.FormValue("title"),
			Author:      r.FormValue("author"),
			Rating:      rating,
			Year:        year,
			Description: r.FormValue("description"),
			Type:        BookType(r.FormValue("type")),
		}

		//books = append(books, book)
		_, err = db.Exec(
			`INSERT INTO books (title, author, rating, year, description, type)
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			book.Title,
			book.Author,
			book.Rating,
			book.Year,
			book.Description,
			book.Type,
		)

		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
