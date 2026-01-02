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

## フロント仕様

### 画面構成
- **ページパス**: `/subscriptions`
- **レイアウト**: 一覧
- **認証**: 必要
- **ナビゲーション**: メニューに表示する

### UI要件

#### 一覧画面 (`/subscriptions`)
- **テーブル表示項目**:
  - サービス名 (name)
  - 金額 (price) - フォーマット: `¥{price}`
  - 次回課金日 (next_billing_date) - フォーマット: YYYY-MM-DD
  - 課金サイクル (billing_cycle) - 表示: "月額" / "年額"
  - 説明 (description) - オプション
- **ソート機能**: 次回課金日順（デフォルト）
- **フィルタ機能**: なし（将来実装可能）
- **ページネーション**: なし（将来実装可能）
- **新規作成ボタン**: あり（`/subscriptions/new` に遷移）
- **行クリック**: 詳細画面に遷移（`/subscriptions/[id]`）
- **アクション**:
  - 編集ボタン: `/subscriptions/[id]/edit` に遷移
  - 削除ボタン: 削除確認モーダル表示後、削除実行

### ルーティング
- 一覧: `/subscriptions`
- 詳細: `/subscriptions/[id]`
- 新規作成: `/subscriptions/new`
- 編集: `/subscriptions/[id]/edit`

