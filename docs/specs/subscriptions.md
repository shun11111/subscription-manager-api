# サブスクリプション管理API

## 概要
サブスクリプションのCRUD操作

## エンドポイント

### GET /subscriptions
- **説明**: サブスクリプション一覧取得
- **認証**: 必要 (Bearer Token)
- **リクエスト**: なし
- **レスポンス**: Subscription[]
- **エラー**: 401 (認証エラー)

### POST /subscriptions
- **説明**: サブスクリプション作成
- **認証**: 必要 (Bearer Token)
- **リクエスト**: CreateSubscriptionRequest
- **レスポンス**: Subscription
- **エラー**: 400 (バリデーションエラー), 401 (認証エラー)

### GET /subscriptions/{id}
- **説明**: サブスクリプション詳細取得
- **認証**: 必要 (Bearer Token)
- **パラメータ**: 
  - `id`: uuid (path)
- **リクエスト**: なし
- **レスポンス**: Subscription
- **エラー**: 401 (認証エラー), 404 (見つからない)

### PUT /subscriptions/{id}
- **説明**: サブスクリプション更新
- **認証**: 必要 (Bearer Token)
- **パラメータ**: 
  - `id`: uuid (path)
- **リクエスト**: UpdateSubscriptionRequest
- **レスポンス**: Subscription
- **エラー**: 400 (バリデーションエラー), 401 (認証エラー), 404 (見つからない)

## データモデル

### Subscription
- `id`: uuid (必須)
- `user_id`: uuid (必須)
- `name`: string (必須)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `next_billing_date`: date (必須)
- `description`: string (オプション)
- `created_at`: datetime (必須)
- `updated_at`: datetime (必須)

### CreateSubscriptionRequest
- `name`: string (必須)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `next_billing_date`: date (必須)
- `description`: string (オプション)

### UpdateSubscriptionRequest
- `name`: string (必須)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `next_billing_date`: date (必須)
- `description`: string (オプション)

