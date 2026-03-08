package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var books []Book

func main() {
	InitDB()

	tmpl, err := template.ParseFiles("static/page.html")
	if err != nil {
		log.Fatal("Template parse error:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		books, err := GetBooks()
		if err != nil {
			log.Println(err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		var readBooks []Book
		var toReadBooks []Book
		var recBooks []Book

		for _, b := range books {

			switch b.Type {
			case Read:
				readBooks = append(readBooks, b)

			case ToRead:
				toReadBooks = append(toReadBooks, b)

			case Rec:
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

		err = AddBook(book)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to save book", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
