# 認証API

## 概要
ユーザー認証・登録に関するAPI

## エンドポイント

### POST /auth/signup
- **説明**: ユーザー登録
- **認証**: 不要
- **リクエスト**: SignUpRequest
- **レスポンス**: SignUpResponse
- **エラー**: 400 (バリデーションエラー)

### POST /auth/login
- **説明**: ログイン
- **認証**: 不要
- **リクエスト**: LoginRequest
- **レスポンス**: LoginResponse
- **エラー**: 401 (認証失敗)

## データモデル

### SignUpRequest
- `email`: string (必須, email形式)
- `password`: string (必須, 8文字以上)
- `name`: string (必須, 1文字以上)

### SignUpResponse
- `user_id`: uuid
- `token`: string (JWTトークン)

### LoginRequest
- `email`: string (必須, email形式)
- `password`: string (必須)

### LoginResponse
- `token`: string (JWTトークン)

