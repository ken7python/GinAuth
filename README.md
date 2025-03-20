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
SECRET_KEY=your_secret_key
```

## JWT Secret Key の生成
JWTのSECRET_KEYをランダムに生成するスクリプトがあります。

### 使い方
```
cd tools
go run generate_secret.go
```
### .envに以下を追加
```
SECRET_KEY=生成されたキー
```

## 実行方法
```
go run main.go database.go user.go
```

## ユーザー認証ページ
### 1.サインアップ
`/signup`
ユーザー登録ページにアクセスできます

### 2.サインイン
`/signup`
サインインページにアクセスできます

### 3.トップページ
`/`
トップページです

## APIの使い方
### 1.ユーザー登録
**リクエスト**
`POST api/register`

```
{
  "username": "testuser",
  "password": "password123"
}
```

**レスポンス**
```
{
    "message":"user created"
}
```
### 2.ログイン
**リクエスト**
`POST api/login`

```
{
  "username": "testuser",
  "password": "password123"
}
```
**レスポンス**
{
    "token": "token..."
}

### 3. プロフィール取得（要認証）
`GET api/profile` （JWTトークンが必要）

**ヘッダー**
```
Authorization: Bearer <your-jwt-token>
```

**レスポンス**
```
{"ID":1,"Username":"testuser"}
```
