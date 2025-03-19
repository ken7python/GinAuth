package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})
	r.GET("/signup", func(c *gin.Context) {
		c.File("templates/signup.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("templates/login.html")
	})

	r.POST("api/register", register)
	r.POST("api/login", login)
	r.GET("api/profile", authMiddleware(), profile)
	r.Run(":8080")
}
