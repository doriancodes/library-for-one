package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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

var books []Book

func main() {

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

		for _, b := range books {
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

		books = append(books, book)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
