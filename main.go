package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var secretKey = "your_secret_key"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Password string `gorm:"unique/not null"`
}

type Claims struct {
	UserID uint `json:"user_id`
	jwt.RegisteredClaims
}

func main() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = database
	db.AutoMigrate(&User{})
	r := gin.Default()
	/*
		r.POST("/register", register)
		r.POST("/login", login)
		r.GET("/users", getUsers)
	*/
	r.Run(":8080")
}
