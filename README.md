# Subscription Manager API

サブスクリプション管理API（仕様駆動・拡張重視）

## 設計コンセプト

- **Schema-First:** `docs/openapi.yaml` を唯一の真実（Single Source of Truth）とし、Goのサーバーコードとフロントの型をここから生成する
- **Clean Architecture:** レイヤ分離による保守性とテスタビリティの向上
- **Standard Layers:** domain, infrastructure, usecase, interface の4層アーキテクチャ

詳細は [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) を参照してください。

## セットアップ

**初めてセットアップする方は、[`docs/SETUP.md`](docs/SETUP.md) を参照してください。**  
Go初心者向けに、環境構築からサーバー起動まで詳しく説明しています。

### クイックスタート（既にGoがインストールされている場合）

```bash
# 1. 依存関係のインストール
go mod download

# 2. 環境変数の設定（既存メンバーから.envファイルをもらう）
# .envファイルをプロジェクトルートに配置

# 3. 環境変数を読み込んでサーバー起動
export $(grep -v '^#' .env | xargs)
go run ./cmd/api/main.go
```

### Dockerを使った起動（オプション）

Dockerを使うと、Goをインストールしなくても起動できます：

```bash
# Dockerイメージをビルド
docker build -t subscription-manager-api .

# コンテナを起動
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="your-secret" \
  subscription-manager-api
```

詳細は [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) を参照してください。

## 📚 ドキュメント

**初めての方は [`docs/README.md`](docs/README.md) から始めてください。**  
ドキュメント一覧と読む順序が説明されています。

**AIを使う開発者の方:**  
各AIツールは自動的に専用ファイルを読み込み、そこから [`AI_INSTRUCTIONS.md`](AI_INSTRUCTIONS.md) を参照します。

- **Cursor**: [`@cursor-instruction.md`](@cursor-instruction.md) → `AI_INSTRUCTIONS.md` を参照
- **GitHub Copilot**: [`.github/copilot-instructions.md`](.github/copilot-instructions.md) または [`.copilot/instructions.md`](.copilot/instructions.md) → `AI_INSTRUCTIONS.md` を参照
- **Claude Code**: [`.claude/instructions.md`](.claude/instructions.md) → `AI_INSTRUCTIONS.md` を参照
- **OpenAI Codex**: [`.codex/instructions.md`](.codex/instructions.md) → `AI_INSTRUCTIONS.md` を参照

**👉 開発指示は [`AI_INSTRUCTIONS.md`](AI_INSTRUCTIONS.md) にまとめられています。**

### 主要ドキュメント

- **[docs/README.md](docs/README.md)** - ドキュメント一覧とガイド
- **[docs/SETUP.md](docs/SETUP.md)** - 環境構築ガイド（Go初心者向け）
- **[docs/ENV_SETUP.md](docs/ENV_SETUP.md)** - 環境変数の設定方法
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - アーキテクチャの説明
- **[docs/DIRECTORY_STRUCTURE.md](docs/DIRECTORY_STRUCTURE.md)** - ディレクトリ構成の説明
- **[docs/WORKFLOW.md](docs/WORKFLOW.md)** - 開発ワークフロー
- **[docs/INTEGRATED_DEVELOPMENT_FLOW.md](docs/INTEGRATED_DEVELOPMENT_FLOW.md)** - **統合開発フロー（フロント+バックエンド）** ⭐ 新画面追加時は必ず参照
- **[docs/API_TESTING.md](docs/API_TESTING.md)** - APIのテスト方法
- **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** - Dockerとデプロイの詳細

## API エンドポイント

### 認証

- `POST /api/auth/signup` - ユーザー登録
- `POST /api/auth/login` - ログイン

### サブスクリプション（認証必要）

- `GET /api/subscriptions` - サブスクリプション一覧取得
- `POST /api/subscriptions` - サブスクリプション作成

## コード生成

このプロジェクトでは、**Schema-First**のアプローチで、Markdown仕様書からOpenAPIを生成し、さらにGoコードを生成しています。

### コード生成の流れ

1. **Markdown仕様書** (`docs/specs/*.md`) を編集
2. **OpenAPI仕様書** (`docs/openapi.yaml`) を生成
3. **Goコード** (`internal/api/openapi.gen.go`) を生成

### コマンド

```bash
# OpenAPI仕様書を生成（Markdown仕様書から）
make gen-openapi

# Goコードを生成（OpenAPIから）
make gen-code

# 両方を実行
make gen-all
```

詳細は [`docs/WORKFLOW.md`](docs/WORKFLOW.md) を参照してください。

## ディレクトリ構成

```
.
├── cmd/
│   └── api/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── domain/              # ドメインモデル（エンティティ）
│   ├── usecase/             # ビジネスロジック層（ユースケース）
│   ├── interface/           # インターフェース層
│   │   └── http/
│   │       ├── handler/     # HTTPハンドラ
│   │       └── middleware/ # HTTPミドルウェア
│   ├── infrastructure/      # インフラストラクチャ層
│   │   ├── persistence/    # データアクセス層（リポジトリ）
│   │   └── config/          # 設定管理
│   └── api/                 # 生成されたOpenAPIコード
├── docs/                    # ドキュメント
│   ├── specs/              # Markdown仕様書
│   └── openapi.yaml        # OpenAPI仕様書
├── migrations/             # データベースマイグレーション
└── pkg/                    # 外部に公開可能なライブラリ
    └── specgen/            # コード生成ライブラリ
```

詳細は [`docs/DIRECTORY_STRUCTURE.md`](docs/DIRECTORY_STRUCTURE.md) を参照してください。

