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
	password = "******"
)

var db *sql.DB

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	router.POST("/login", LoginHandler)

	router.Run("localhost:8080")
}

func getProducts(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	// Remove "Bearer " prefix
	tokenString := authHeader[len("Bearer "):]

	// Verify token
	if err := VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Query the DB
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
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

func LoginHandler(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Query the user by username
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", u.Username).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if password matches
	if u.Password != dbPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token
	tokenString, err := CreateToken(u.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Success
	c.String(http.StatusOK, tokenString)
}
