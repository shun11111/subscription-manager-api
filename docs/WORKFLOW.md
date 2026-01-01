# 開発ワークフロー

## 📖 このドキュメントについて

このドキュメントでは、新しい画面/APIを追加する際のワークフローを説明します。

**対象読者:**
- 新機能を追加する方
- 開発の流れを理解したい方

**読む前に:**
- 環境構築が完了していること（[`SETUP.md`](./SETUP.md) を参照）
- アーキテクチャを理解していること（[`ARCHITECTURE.md`](./ARCHITECTURE.md) を参照）

**関連ドキュメント:**
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの説明
- [`DIRECTORY_STRUCTURE.md`](./DIRECTORY_STRUCTURE.md) - ディレクトリ構成
- [`API_TESTING.md`](./API_TESTING.md) - APIのテスト方法

---

## 概要

このプロジェクトでは、**Markdown形式の仕様書**をベースにAPIを管理します。

## 新しい画面を追加する手順

### 1. 仕様書を作成

`docs/specs/{画面名}.md` を作成し、テンプレートに従って記述します。

例：`docs/specs/billing.md` を作成

```markdown
# 請求管理API

## 概要
請求履歴の管理

## エンドポイント

### GET /billing
- **説明**: 請求履歴一覧取得
- **認証**: 必要
- **リクエスト**: なし
- **レスポンス**: Billing[]
- **エラー**: 401 (認証エラー)

## データモデル

### Billing
- `id`: uuid (必須)
- `subscription_id`: uuid (必須)
- `amount`: float (必須)
- `created_at`: datetime (必須)
```

### 2. OpenAPIを生成

```bash
make gen-openapi
```

これで `docs/openapi.yaml` が自動更新されます。

### 3. Goコードを生成

```bash
make gen-code
```

これで `internal/api/openapi.gen.go` が生成されます。

### 4. Infrastructure/Persistence層を実装

`internal/infrastructure/persistence/{画面名}_repository.go` を作成し、SQLを実装します。

### 5. Usecase層を実装

`internal/usecase/{画面名}_service.go` を作成し、ビジネスロジックを実装します。

### 6. Interface/HTTP/Handler層を実装

`internal/interface/http/handler/{画面名}_handler.go` を作成し、HTTPハンドラを実装します。

### 7. ルーティングを追加

`cmd/api/main.go` にルーティングを追加します。

### 8. テストケースの生成

`docs/tests/{画面名}_test.sh` または `docs/tests/{画面名}_test.md` を作成し、テストケースを記述します。

- 認証が必要な場合: まずログインしてトークンを取得する手順
- 各エンドポイント（GET, POST, PUT, DELETE）のcurlコマンド例
- 成功時のレスポンス例
- エラーケース（401, 400, 404など）のテスト

詳細は [`API_TESTING.md`](./API_TESTING.md) を参照してください。

## 既存画面を編集する手順

1. `docs/specs/{画面名}.md` を編集
2. `make gen-openapi` でOpenAPIを更新
3. `make gen-code` でGoコードを更新
4. Infrastructure/Usecase/Handler層を必要に応じて修正

## メリット

- **YAMLを直接書かなくて良い**: Markdownの方が読み書きしやすい
- **仕様書が常に最新**: コード生成の元になるので、仕様書と実装が乖離しない
- **画面ごとにファイル分割**: 管理しやすい
- **テンプレートで統一**: 新しい画面も同じ形式で書ける

