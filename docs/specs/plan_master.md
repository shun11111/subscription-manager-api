# プラン・価格マスター管理画面

## 概要
プランと価格のマスターデータを管理するAPI

## エンドポイント

### GET /plans
- **説明**: プラン一覧取得
- **認証**: 必要 (Bearer Token)
- **リクエスト**: なし
- **レスポンス**: Plan[]
- **エラー**: 401 (認証エラー)

### POST /plans
- **説明**: プラン作成
- **認証**: 必要 (Bearer Token)
- **リクエスト**: CreatePlanRequest
- **レスポンス**: Plan
- **エラー**: 400 (バリデーションエラー), 401 (認証エラー)

### GET /plans/{id}
- **説明**: プラン詳細取得
- **認証**: 必要 (Bearer Token)
- **パラメータ**: 
  - `id`: uuid (path)
- **リクエスト**: なし
- **レスポンス**: Plan
- **エラー**: 401 (認証エラー), 404 (見つからない)

### PUT /plans/{id}
- **説明**: プラン更新
- **認証**: 必要 (Bearer Token)
- **パラメータ**: 
  - `id`: uuid (path)
- **リクエスト**: UpdatePlanRequest
- **レスポンス**: Plan
- **エラー**: 400 (バリデーションエラー), 401 (認証エラー), 404 (見つからない)

### DELETE /plans/{id}
- **説明**: プラン削除
- **認証**: 必要 (Bearer Token)
- **パラメータ**: 
  - `id`: uuid (path)
- **リクエスト**: なし
- **レスポンス**: なし
- **エラー**: 401 (認証エラー), 404 (見つからない)

## データモデル

### Plan
- `id`: uuid (必須)
- `name`: string (必須)
- `description`: string (オプション)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `features`: string[] (オプション)
- `is_active`: boolean (必須)
- `created_at`: datetime (必須)
- `updated_at`: datetime (必須)

### CreatePlanRequest
- `name`: string (必須)
- `description`: string (オプション)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `features`: string[] (オプション)
- `is_active`: boolean (必須)

### UpdatePlanRequest
- `name`: string (必須)
- `description`: string (オプション)
- `price`: float (必須)
- `billing_cycle`: enum (必須, "monthly" | "yearly")
- `features`: string[] (オプション)
- `is_active`: boolean (必須)

