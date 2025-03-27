package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Database connection string
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Eros2724@1"
)

var db *sql.DB

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func main() {
	// Connect to the database
	psqlInfo := "host=localhost port=5432 user=" + user + " password=" + password + " sslmode=disable"
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	db = dbConn
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	println("Connected to PostgreSQL successfully!")

	// Setup routes
	router := gin.Default()
	router.GET("/products", getProducts)

	router.Run("localhost:8080")
}

func getProducts(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Stock, &prod.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data scan error"})
			return
		}
		products = append(products, prod)
	}

	c.JSON(http.StatusOK, products)
}
