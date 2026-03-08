package main

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
