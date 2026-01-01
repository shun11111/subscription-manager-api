# API仕様書

このディレクトリには、**画面ごとのAPI仕様をMarkdown形式で管理**します。

## 使い方

1. **新しい画面を追加するとき**
   - `docs/specs/{画面名}.md` を作成
   - テンプレートに従ってAPI仕様を記述
   - `make gen-openapi` を実行して `docs/openapi.yaml` を自動生成

2. **既存の画面を編集するとき**
   - `docs/specs/{画面名}.md` を編集
   - `make gen-openapi` を実行

## ファイル構成

- `auth.md` - 認証関連API
- `subscriptions.md` - サブスクリプション管理API
- `template.md` - 新しい画面を追加するときのテンプレート

## 仕様書の形式

各仕様書は以下の形式で記述します：

```markdown
# {画面名}

## 概要
画面の説明

## エンドポイント

### GET /resource
- **説明**: 一覧取得
- **認証**: 必要
- **リクエスト**: なし
- **レスポンス**: Resource[]

### POST /resource
- **説明**: 作成
- **認証**: 必要
- **リクエスト**: CreateResourceRequest
- **レスポンス**: Resource

## データモデル

### Resource
- `id`: uuid (必須)
- `name`: string (必須)
- `created_at`: datetime (必須)

### CreateResourceRequest
- `name`: string (必須)
- `description`: string (オプション)
```

