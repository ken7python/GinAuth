package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const secretLength = 32

func main() {
	secret := make([]byte, secretLength)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Println("エラー: 秘密鍵を生成できませんでした")
		return
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)
	fmt.Println(encodedSecret)
}
