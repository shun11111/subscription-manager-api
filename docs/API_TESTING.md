# APIテストガイド

## 📖 このドキュメントについて

このドキュメントでは、APIをテストする方法を説明します。

**対象読者:**
- APIをテストしたい方
- 動作確認を行いたい方

**読む前に:**
- 環境構築が完了していること（[`SETUP.md`](./SETUP.md) を参照）
- 環境変数が設定されていること（[`ENV_SETUP.md`](./ENV_SETUP.md) を参照）
- サーバーが起動していること

**関連ドキュメント:**
- [`SETUP.md`](./SETUP.md) - 環境構築ガイド
- [`ENV_SETUP.md`](./ENV_SETUP.md) - 環境変数の設定
- [`WORKFLOW.md`](./WORKFLOW.md) - 開発ワークフロー

---

このドキュメントでは、Subscription Manager APIをテストする方法を説明します。

## 前提条件

- サーバーが起動していること（`go run ./cmd/api/main.go`）
- データベース（Supabase）にマイグレーションが実行されていること
- 環境変数（`.env`）が正しく設定されていること

## テスト方法

### 1. curlコマンドを使う方法（推奨）

#### 1-1. ユーザー登録（サインアップ）

```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

**成功時のレスポンス例:**
```json
{
  "user_id": "1e78e0dc-8663-4347-805c-fed2473b7ac1",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンスから`token`をコピーして、以下の環境変数に設定:**
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### 1-2. ログイン

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**成功時のレスポンス例:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 1-3. プラン一覧取得

```bash
curl -X GET http://localhost:8080/api/plans \
  -H "Authorization: Bearer $TOKEN"
```

#### 1-4. プラン作成

```bash
curl -X POST http://localhost:8080/api/plans \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Basic Plan",
    "description": "ベーシックプラン",
    "price": 1000,
    "billing_cycle": "monthly",
    "features": ["feature-a", "feature-b"],
    "is_active": true
  }'
```

#### 1-5. プラン詳細取得

```bash
# プランIDを指定（例: c2684fe8-6cf8-4ca0-ae7d-ae4876e8d4c9）
curl -X GET http://localhost:8080/api/plans/c2684fe8-6cf8-4ca0-ae7d-ae4876e8d4c9 \
  -H "Authorization: Bearer $TOKEN"
```

#### 1-6. プラン更新

```bash
curl -X PUT http://localhost:8080/api/plans/c2684fe8-6cf8-4ca0-ae7d-ae4876e8d4c9 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Basic Plan Updated",
    "description": "更新されたベーシックプラン",
    "price": 1200,
    "billing_cycle": "monthly",
    "features": ["feature-a", "feature-b", "feature-c"],
    "is_active": true
  }'
```

#### 1-7. プラン削除

```bash
curl -X DELETE http://localhost:8080/api/plans/c2684fe8-6cf8-4ca0-ae7d-ae4876e8d4c9 \
  -H "Authorization: Bearer $TOKEN"
```

#### 1-8. サブスクリプション一覧取得

```bash
curl -X GET http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer $TOKEN"
```

#### 1-9. サブスクリプション作成

```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix",
    "price": 990,
    "billing_cycle": "monthly",
    "next_billing_date": "2026-02-01",
    "description": "動画配信サービス"
  }'
```

#### 1-10. サブスクリプション詳細取得

```bash
curl -X GET http://localhost:8080/api/subscriptions/ba9b37ea-d29d-421b-b613-8863d2ab5dbd \
  -H "Authorization: Bearer $TOKEN"
```

#### 1-11. サブスクリプション更新

```bash
curl -X PUT http://localhost:8080/api/subscriptions/ba9b37ea-d29d-421b-b613-8863d2ab5dbd \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix Premium",
    "price": 1490,
    "billing_cycle": "monthly",
    "next_billing_date": "2026-02-01",
    "description": "プレミアムプラン"
  }'
```

### 2. HTTPieを使う方法（より読みやすい）

HTTPieをインストール:
```bash
# macOS
brew install httpie

# または
pip install httpie
```

#### 使用例

```bash
# サインアップ
http POST http://localhost:8080/api/auth/signup \
  email=test@example.com \
  password=password123 \
  name="Test User"

# ログイン
http POST http://localhost:8080/api/auth/login \
  email=test@example.com \
  password=password123

# プラン一覧取得
http GET http://localhost:8080/api/plans \
  Authorization:"Bearer $TOKEN"

# プラン作成
http POST http://localhost:8080/api/plans \
  Authorization:"Bearer $TOKEN" \
  name="Basic Plan" \
  description="ベーシックプラン" \
  price:=1000 \
  billing_cycle=monthly \
  features:='["feature-a","feature-b"]' \
  is_active:=true
```

### 3. Postman / Insomniaを使う方法

#### セットアップ

1. **コレクションを作成**
   - Postman/Insomniaで新しいコレクションを作成
   - ベースURL: `http://localhost:8080/api`

2. **環境変数を設定**
   - `base_url`: `http://localhost:8080/api`
   - `token`: （ログイン後に取得したトークン）

3. **リクエスト例**

   **認証 > サインアップ**
   - Method: `POST`
   - URL: `{{base_url}}/auth/signup`
   - Body (JSON):
     ```json
     {
       "email": "test@example.com",
       "password": "password123",
       "name": "Test User"
     }
     ```

   **認証 > ログイン**
   - Method: `POST`
   - URL: `{{base_url}}/auth/login`
   - Body (JSON):
     ```json
     {
       "email": "test@example.com",
       "password": "password123"
     }
     ```

   **プラン > 一覧取得**
   - Method: `GET`
   - URL: `{{base_url}}/plans`
   - Headers:
     - `Authorization`: `Bearer {{token}}`

   **プラン > 作成**
   - Method: `POST`
   - URL: `{{base_url}}/plans`
   - Headers:
     - `Authorization`: `Bearer {{token}}`
     - `Content-Type`: `application/json`
   - Body (JSON):
     ```json
     {
       "name": "Basic Plan",
       "description": "ベーシックプラン",
       "price": 1000,
       "billing_cycle": "monthly",
       "features": ["feature-a", "feature-b"],
       "is_active": true
     }
     ```

### 4. ブラウザでテスト（GETリクエストのみ）

認証が必要なエンドポイントはブラウザから直接アクセスできませんが、開発用に一時的に認証を外すか、ブラウザ拡張機能（ModHeaderなど）でトークンを設定する必要があります。

## エラーレスポンス

### 401 Unauthorized
```json
{
  "message": "Authorization header is required"
}
```
**対処法**: `Authorization: Bearer <token>` ヘッダーを追加

### 400 Bad Request
```json
{
  "message": "Invalid request"
}
```
**対処法**: リクエストボディの形式を確認

### 404 Not Found
```json
{
  "message": "not found"
}
```
**対処法**: 指定したIDが存在するか確認

### 500 Internal Server Error
```json
{
  "message": "error message"
}
```
**対処法**: サーバーログを確認

## テストスクリプト例

`test_api.sh` を作成:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api"

# サインアップ
echo "=== サインアップ ==="
SIGNUP_RESPONSE=$(curl -s -X POST $BASE_URL/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}')
echo $SIGNUP_RESPONSE | jq .

# トークンを取得
TOKEN=$(echo $SIGNUP_RESPONSE | jq -r '.token')
export TOKEN

# ログイン
echo -e "\n=== ログイン ==="
curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq .

# プラン作成
echo -e "\n=== プラン作成 ==="
curl -s -X POST $BASE_URL/plans \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Basic Plan",
    "description": "ベーシックプラン",
    "price": 1000,
    "billing_cycle": "monthly",
    "features": ["feature-a"],
    "is_active": true
  }' | jq .

# プラン一覧取得
echo -e "\n=== プラン一覧取得 ==="
curl -s -X GET $BASE_URL/plans \
  -H "Authorization: Bearer $TOKEN" | jq .
```

実行:
```bash
chmod +x test_api.sh
./test_api.sh
```

## トラブルシューティング

### サーバーが起動しない
- ポート8080が使用中: `lsof -ti:8080 | xargs kill -9`
- 環境変数が設定されていない: `.env`ファイルを確認

### 認証エラー
- トークンが期限切れ: 再度ログインしてトークンを取得
- トークンの形式が間違っている: `Bearer ` の後にスペースが必要

### データベースエラー
- マイグレーションが実行されていない: Supabase SQL Editorで実行
- 接続文字列が間違っている: `.env`の`DATABASE_URL`を確認

