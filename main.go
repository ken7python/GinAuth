package main

import (
	"github.com/gin-gonic/gin"
)

/*
var db *gorm.DB

var secretKey = "your_secret_key"
*/

func main() {
	InitDB()
	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/profile", authMiddleware(), profile)
	r.Run(":8080")
}
