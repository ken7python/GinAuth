package main

import (
	"os"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var db *gorm.DB
var secretKey = "your_secret_key"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique/not null"`
	Password string `gorm:"unique/not null"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	//println(dsn)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = database
	db.AutoMigrate(&User{})
	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	// r.GET("/users", getUsers)
	r.Run(":8080")
}

func register(c *gin.Context) {
	var req struct {
		Username string `json:"username""`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	db.Create(&User{Username: req.Username, Password: string(hashedPassword)})

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token := createToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createToken(userID uint) string {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString

}
