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
	"github.com/dgrijalva/jwt-go"
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

// Global rate limiting map (intentionally insecure)
var requestCounts = make(map[string]int)

// CustomClaims represents the JWT claims structure
type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// parseToken parses a JWT token string and returns the token object
func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}

// parseTokenWithClaims parses a JWT token string with custom claims
func parseTokenWithClaims(tokenString string) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	return token, claims, err
}

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
			// Create JWT token with custom claims
			claims := &CustomClaims{
				Username: username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
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

	// Vulnerable product management endpoints
	r.POST("/api/products", func(c *gin.Context) {
		name := c.PostForm("name")
		price := c.PostForm("price")
		description := c.PostForm("description")

		// SQL Injection vulnerability in product creation
		query := fmt.Sprintf("INSERT INTO products (name, price, description) VALUES ('%s', %s, '%s')", 
			name, price, description)
		_, err := db.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	})

	// Vulnerable file upload endpoint
	r.POST("/api/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		defer file.Close()

		// Path traversal vulnerability
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

	// Vulnerable comment system with stored XSS
	r.POST("/api/comments", func(c *gin.Context) {
		productID := c.PostForm("product_id")
		comment := c.PostForm("comment")
		username := c.PostForm("username")

		// Stored XSS vulnerability
		query := fmt.Sprintf("INSERT INTO comments (product_id, username, comment) VALUES (%s, '%s', '%s')",
			productID, username, comment)
		_, err := db.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
	})

	// Vulnerable rate limiting implementation
	r.Use(func(c *gin.Context) {
		ip := c.ClientIP()
		requestCounts[ip]++
		
		// Rate limiting bypass vulnerability - can be bypassed by changing IP in headers
		if requestCounts[ip] > 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Reset counter after 1 minute (intentionally insecure)
		go func() {
			time.Sleep(time.Minute)
			requestCounts[ip] = 0
		}()
		
		c.Next()
	})

	// Add a new protected endpoint that uses token validation
	r.GET("/api/protected", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract the token from the Authorization header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse the token with claims
		token, claims, err := parseTokenWithClaims(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Protected endpoint accessed successfully",
			"user":    claims.Username,
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
} 