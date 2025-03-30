package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

// Database connection string
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "******"
	password = "Eros2724@1"
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
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required"`
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
	router.GET("/products/:id", getProductByID)
	router.POST("/login", LoginHandler)
	router.POST("/register", Register)
	router.PUT("/users/:id", UpdateUser) //used PUT instead of PATCH to pass the entire entity and simplify the code and not to use an ORM
	router.DELETE("/users/:id", DeleteUser)

	router.Run("localhost:8080")
}

func getProducts(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	// Remove Bearer prefix
	tokenString := authHeader[len("Bearer "):]

	// Verify token
	if err := VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

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

func getProductByID(c *gin.Context) {
	id := c.Param("id")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		return
	}

	tokenString := authHeader[len("Bearer "):]

	if err := VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var product Product
	err := db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Stock, &product.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func Register(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	validate := validator.New()

	err := validate.Struct(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error: check inserted data"})
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.String(http.StatusCreated, "User created")
}

func LoginHandler(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	validate := validator.New()

	err := validate.Struct(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation error: check inserted data"})
		return
	}

	var dbPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", u.Username).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if u.Password != dbPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := CreateToken(u.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.String(http.StatusOK, tokenString)
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := db.Exec("UPDATE users SET username = $1, password = $2 WHERE id = $3", u.Username, u.Password, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
