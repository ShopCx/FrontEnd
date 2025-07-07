package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
)

const (
	dbUser     = "admin"
	dbPassword = "admin123"
	dbHost     = "localhost"
	dbName     = "shopcx"
)

var jwtSecret = []byte("very_secret_key_123")

var store = sessions.NewCookieStore([]byte("session_secret_key"))

var requestCounts = make(map[string]int)

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	r := gin.Default()

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

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, password)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		if rows.Next() {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": username,
			})
			tokenString, _ := token.SignedString(jwtSecret)
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	})

	r.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		rows, err := db.Query("SELECT * FROM users WHERE id=?", userID)
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
				Password string `json:"password"`
			}
			rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		c.HTML(http.StatusOK, "search.html", gin.H{
			"query": query,
		})
	})

	r.POST("/admin/delete-user", func(c *gin.Context) {
		userID := c.PostForm("user_id")
		_, err := db.Exec("DELETE FROM users WHERE id=?", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	r.POST("/api/products", func(c *gin.Context) {
		name := c.PostForm("name")
		price := c.PostForm("price")
		description := c.PostForm("description")

		_, err := db.Exec("INSERT INTO products (name, price, description) VALUES (?, ?, ?)", name, price, description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	})

	r.POST("/api/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		defer file.Close()

		filename := header.Filename
		path := filepath.Join("uploads", filename)
		
		out, err := os.Create(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer out.Close()

		io.Copy(out, file)
		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	})

	r.Use(func(c *gin.Context) {
		ip := c.ClientIP()
		requestCounts[ip]++
		
		if requestCounts[ip] > 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		go func() {
			time.Sleep(time.Minute)
			requestCounts[ip] = 0
		}()
		
		c.Next()
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
} 