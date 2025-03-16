# GinAuth - Golangで作るJWT認証API

Gin + GORM + MySQLを使ったシンプルなユーザー認証システムです。
JWTを使った認証機能を提供し、安全にユーザー管理を行えます。

## 使用技術
 - Golang 1.24.1
 - Gin (Webフレームワーク)
 - GORM (ORM)
 - MySQL (データベース)
 - JWT (認証)

## 環境変数の設定
.envを以下のようにしてください
```
DB_USER=username
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=yourdatabase
```

## APIの使い方
### 1.ユーザー登録
リクエスト
POST /register

```
{
  "username": "testuser",
  "password": "password123"
}
```

レスポンス
```
{
    "message":"user created"
}
```
### 2.ログイン
リクエスト
POST /login

```
{
  "username": "testuser",
  "password": "password123"
}
```
レスポンス
{
    "token": "token..."
}