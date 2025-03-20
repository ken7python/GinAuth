package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const secretLength = 32

func main() {
	encodedSecret := generate_secret()
	if encodedSecret == "" {
		return
	}
	fmt.Println(encodedSecret)
}

func generate_secret() string {
	secret := make([]byte, secretLength)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Println("エラー: 秘密鍵を生成できませんでした")
		return ""
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)
	return encodedSecret
}
