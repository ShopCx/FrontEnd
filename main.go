package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
)

// Hardcoded database credentials (intentionally insecure)
const (
	dbUser     = "admin"
	dbPassword = "admin123"
	dbHost     = "localhost"
	dbName     = "shopcx"
)

// Hardcoded JWT secret (intentionally insecure)
var jwtSecret = []byte("very_secret_key_123")

// Global session store with weak secret (intentionally insecure)
var store = sessions.NewCookieStore([]byte("session_secret_key"))

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	r := gin.Default()

	// CORS middleware with overly permissive configuration (intentionally insecure)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Vulnerable login endpoint with SQL injection
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// SQL Injection vulnerability (intentionally insecure)
		query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, password)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		if rows.Next() {
			// Create JWT token with weak algorithm (intentionally insecure)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": username,
			})
			tokenString, _ := token.SignedString(jwtSecret)
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	})

	// Vulnerable user profile endpoint with IDOR
	r.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		// IDOR vulnerability (intentionally insecure)
		query := fmt.Sprintf("SELECT * FROM users WHERE id=%s", userID)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		if rows.Next() {
			var user struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
				Password string `json:"password"` // Exposing password (intentionally insecure)
			}
			rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	})

	// Vulnerable search endpoint with XSS
	// Intentionally undocumented in Swagger: Internal search functionality
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		// XSS vulnerability (intentionally insecure)
		c.HTML(http.StatusOK, "search.html", gin.H{
			"query": query,
		})
	})

	// Undocumented admin endpoint (intentionally hidden)
	r.POST("/admin/delete-user", func(c *gin.Context) {
		userID := c.PostForm("user_id")
		// No authentication check (intentionally insecure)
		query := fmt.Sprintf("DELETE FROM users WHERE id=%s", userID)
		_, err := db.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
} 