package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"crypto/rand"
	"encoding/base64"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UUID     string `gorm:"unique/not null"`
	Username string `gorm:"unique/not null"`
	Password string `gorm:"unique/not null"`
}

var secretKey string

func generate_secret() string {
	const secretLength = 32
	secret := make([]byte, secretLength)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Println("エラー: 秘密鍵を生成できませんでした")
		return ""
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)
	return encodedSecret
}

func GetSecretKey() string {
	if secretKey != "" {
		return secretKey
	}

	if err := godotenv.Load(); err == nil {
		secretKey = os.Getenv("SECRET_KEY")
	}

	if secretKey == "" {
		secretKey = generate_secret()
		if secretKey == "" {
			panic("Error: Failed to generate a secret key")
		}
	}
	return secretKey
}

func parseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecretKey()), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func register(c *gin.Context) {
	println("/register")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの解析に失敗しました"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ユーザー名は既に使用されています"})
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUIDの生成に失敗しました"})
		return
	}
	fmt.Printf("uuid: %s\n", uuid.String())

	db.Create(&User{Username: req.Username, Password: string(hashedPassword), UUID: uuid.String()})

	c.JSON(http.StatusOK, gin.H{"message": "ユーザーを作成しました"})
}

func login(c *gin.Context) {
	println("/login")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストの解析に失敗しました"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー名またはパスワードが正しくありません"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー名またはパスワードが正しくありません"})
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
	tokenString, _ := token.SignedString([]byte(GetSecretKey()))

	return tokenString

}

func authMiddleware() gin.HandlerFunc {
	println("authMiddleware")
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			println("tokenString == ''")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
			c.Abort()
			return
		}
		// Remove "Bearer " prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := parseToken(tokenString)
		if err != nil {
			println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func profile(c *gin.Context) {
	println("/profile")
	userID, exists := c.Get("userID")
	if !exists {
		println("!exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ID": user.ID, "username": user.Username})
}
