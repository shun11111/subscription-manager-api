# ディレクトリ構成の説明

## 📖 このドキュメントについて

このドキュメントでは、プロジェクトのディレクトリ構成と各ディレクトリの役割を説明します。

**対象読者:**
- プロジェクト構造を理解したい方
- 新しいファイルを追加する方
- Clean Architectureを学びたい方

**読む前に:**
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) を読むことをおすすめします（アーキテクチャの理解）

**関連ドキュメント:**
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの説明
- [`WORKFLOW.md`](./WORKFLOW.md) - 開発ワークフロー

---

## 現在の構成

```
.
├── cmd/                    # アプリケーションのエントリーポイント
│   ├── api/               # APIサーバー
│   └── specgen/           # コード生成ツール
├── internal/              # プライベートなアプリケーションコード
│   ├── domain/            # ドメインモデル（エンティティ）
│   ├── usecase/           # ビジネスロジック層（ユースケース）
│   ├── interface/         # インターフェース層
│   │   └── http/
│   │       ├── handler/   # HTTPハンドラ層
│   │       └── middleware/ # HTTPミドルウェア
│   ├── infrastructure/    # インフラストラクチャ層
│   │   ├── persistence/   # データアクセス層（リポジトリ）
│   │   └── config/        # 設定管理
│   └── api/               # 生成されたOpenAPIコード
├── pkg/                   # 外部に公開可能なライブラリコード
│   └── specgen/           # コード生成ライブラリ
├── docs/                  # ドキュメント
│   ├── specs/             # Markdown仕様書
│   └── openapi.yaml       # OpenAPI仕様
├── migrations/            # データベースマイグレーション
└── Makefile              # ビルド・生成コマンド
```

## Clean Architectureの観点からの評価

### ✅ 良い点

1. **`internal/` と `pkg/` の分離**
   - `internal/` はプライベート（外部からインポート不可）
   - `pkg/` は再利用可能なライブラリ
   - Goの標準的な慣習に従っている

2. **レイヤの分離**
   - `domain` → `infrastructure` → `usecase` → `interface` の依存関係が明確
   - Clean Architectureの原則に従った構成

3. **`cmd/` の使用**
   - エントリーポイントを分離（Standard Go Project Layout準拠）

### 🔄 改善できる点

#### 1. 命名の改善

現在の命名は機能的ですが、Clean Architectureの用語に合わせるとより明確になります：

**現状:**
- `handler` → HTTPハンドラ（インターフェース層）
- `service` → ビジネスロジック（ユースケース層）
- `repository` → データアクセス（インフラ層）

**Clean Architectureの用語に合わせる場合:**
- `handler` → `interface` または `adapter/infrastructure/http`
- `service` → `usecase` または `application`
- `repository` → `adapter/infrastructure/persistence` または `gateway`

ただし、**Goコミュニティでは現在の命名も一般的**なので、無理に変更する必要はありません。

#### 2. ディレクトリ構造の改善案

よりClean Architectureに近づける場合：

```
.
├── cmd/
│   └── api/
├── internal/
│   ├── domain/              # エンティティ（変更なし）
│   ├── usecase/             # ユースケース（service → usecase）
│   │   ├── auth/
│   │   ├── subscription/
│   │   └── plan/
│   ├── interface/           # インターフェース層（handler → interface）
│   │   └── http/
│   │       ├── handler/
│   │       └── middleware/
│   └── infrastructure/      # インフラ層
│       ├── persistence/     # repository → infrastructure/persistence
│       └── config/          # pkg/config → infrastructure/config
└── pkg/                     # 共有ライブラリ（変更なし）
```

## 現在の構成の評価

### ✅ Clean Architecture準拠

現在の構成は、Clean Architectureの原則に完全に準拠しています：

1. **レイヤの分離**
   - `domain`（最内層）→ `infrastructure` → `usecase` → `interface`（最外層）
   - 依存関係の向きが正しい（内側から外側へ）

2. **明確な責務分離**
   - `domain`: ビジネスエンティティ
   - `usecase`: ビジネスロジック
   - `interface`: 外部インターフェース（HTTP）
   - `infrastructure`: 技術的詳細（DB、設定）

3. **拡張性**
   - 将来的にgRPCなどのインターフェースを追加する場合、`internal/interface/grpc/` を追加するだけで対応可能
   - データベースを変更する場合、`internal/infrastructure/persistence/` を差し替えるだけで対応可能

### 構成の利点

1. **Clean Architectureの用語に一致**
   - `usecase`, `interface`, `infrastructure` は標準的な用語
   - チームメンバーが理解しやすい

2. **実用性**
   - 依存関係が明確で、テストしやすい
   - 各レイヤの責務が明確

3. **保守性**
   - 変更の影響範囲が限定的
   - コードの可読性が高い

## 参考：Standard Go Project Layout

Goの標準的なプロジェクト構成（[golang-standards/project-layout](https://github.com/golang-standards/project-layout)）では：

- `cmd/` - メインアプリケーション
- `internal/` - プライベートなアプリケーションコード
- `pkg/` - 外部に公開可能なライブラリ
- `api/` - API定義（OpenAPIなど）

現在の構成は、この標準に準拠しています。

