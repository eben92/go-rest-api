package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// book represents a book with its ID, title, author, and quantity.
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Golang pointers", Author: "Mr. Golang", Quantity: 2},
	{ID: "2", Title: "Goroutines", Author: "Mr. Goroutine", Quantity: 20},
	{ID: "3", Title: "Golang routers", Author: "Mr. Router", Quantity: 30},
	{ID: "4", Title: "Golang concurrency", Author: "Mr. Currency", Quantity: 40},
}

// getBooks returns a list of all books.
// It takes a pointer to a gin.Context object as its only parameter.
// It uses the IndentedJSON method of the gin.Context object to send an HTTP response with the list of books in JSON format.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// createBook creates a new book and appends it to the books slice.
// It expects a JSON payload in the request body with the following format:
//
//	{
//	  "title": "string",
//	  "author": "string",
//	  "quantity": "int",
//	  "id": "string"
//	}
//
// It returns the newly created book as a JSON response with status code 201 (Created).
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// bookById handles GET requests for a single book by ID.
func bookById(c *gin.Context) {
	id := c.Param("id")

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// getBookById returns a pointer to a book and an error. It takes a string id as input.
// It searches for a book in the books slice with the given id and returns a pointer to the book if found.
// If the book is not found, it returns nil and an error.
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// checkoutBook is a handler function that checks out a book by its ID.
// It retrieves the book by ID, decrements its quantity by 1, and returns the updated book.
// If the book is not found or its quantity is 0, it returns an error message.
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter 'id' "})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book is not available at the moment, check in again later"})
		return
	}

	book.Quantity -= 1

	c.IndentedJSON(http.StatusOK, gin.H{"message": "success", "data": book})
}

// returnBook returns a book by its ID and increments its quantity by 1.
// If the book is not found, it returns a 404 status code.
// If the 'id' query parameter is missing, it returns a 400 status code.
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameter 'id'"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	// Creates a new default Gin router instance with logging and recovery middleware.
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:3001")

}
