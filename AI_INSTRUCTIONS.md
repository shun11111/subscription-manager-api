# subscription-manager-api 開発インストラクション

> **📌 このファイルについて**  
> このファイルは **AI（Cursor、GitHub Copilot、Claude Code、Codexなど）が読むための指示書** です。  
> 新しい画面を追加する際、AIはこのファイルを参照して自動的に実装を行います。  
> 人間が読む場合は、[`docs/WORKFLOW.md`](docs/WORKFLOW.md) や [`docs/README.md`](docs/README.md) を参照してください。
>
> **各AIツールからの参照:**
> - Cursor: [`@cursor-instruction.md`](@cursor-instruction.md) → このファイルを参照
> - GitHub Copilot: [`.github/copilot-instructions.md`](.github/copilot-instructions.md) または [`.copilot/instructions.md`](.copilot/instructions.md) → このファイルを参照
> - Claude Code: [`.claude/instructions.md`](.claude/instructions.md) → このファイルを参照
> - OpenAI Codex: [`.codex/instructions.md`](.codex/instructions.md) → このファイルを参照

---

## 1. このプロジェクトの前提（AI向け）

- **Schema-First:**  
  - API仕様のソースは `docs/openapi.yaml`。  
  - ただし、実際に編集するのは `docs/specs/*.md`（Markdown仕様書）。  
  - `make gen-openapi` で Markdown → `openapi.yaml` を生成する。
- **コード生成:**  
  - `internal/api/oapi_gen.go` の `//go:generate`（`oapi-codegen`）で  
    `docs/openapi.yaml` → `internal/api/openapi.gen.go`（types + server + echo）を生成。
  - `make gen-code` で `go generate ./...` を叩く。
- **レイヤ構成（Go側）:**  
  - `internal/interface/http/handler/*` … HTTPハンドラ（生成された型を使う）  
  - `internal/usecase/*` … ビジネスロジック（ユースケース）  
  - `internal/infrastructure/persistence/*` … DBアクセス（Genericsベースの共通Repositoryあり）  
  - `internal/interface/http/middleware/auth.go` … JWT認証ミドルウェア（`user_id` を context に詰める）

すでに基盤は出来ているので、**AIは「新しい画面追加」だけを意識すればよい**。

---

## 2. ユーザーがAIに伝える一言（テンプレ）

**👤 人間向け:** 新しい画面を追加する際は、以下のテンプレートを使ってAIに指示してください。

**🤖 AI向け:** ユーザーから以下のような指示が来た場合、このファイルの手順に従って実装してください。

### テンプレート

> 「`{画面名}` という管理画面を追加したいです。  
> `docs/specs` 以下のMarkdown仕様書を起点に、OpenAPI → Goコード → infrastructure/usecase/handler → routing まで、  
> このリポジトリのルールに従って一通り実装してください。  
> また、テストケースも自動生成してください。」

**🤖 AI向け:** このファイルと `docs/WORKFLOW.md` / `docs/specs/*.md` / `docs/API_TESTING.md` を前提に動くこと。

---

## 3. 新しい画面を追加するときのAIの手順

**🤖 AI向け:** 以下の手順を**可能な限り自動でツールを叩きながら**実行すること。

**👤 人間向け:** このセクションはAIが自動的に実行する手順です。手動で実装する場合は [`docs/WORKFLOW.md`](docs/WORKFLOW.md) を参照してください。

### 3-1. 仕様書（Markdown）の作成/更新

1. `docs/specs/template.md` をベースに、`docs/specs/{画面名}.md` を作成する。  
   - 既存画面の場合は、該当する `docs/specs/*.md` を編集。
2. エンドポイントとデータモデルを記述する。  
   - 形式は `auth.md` / `subscriptions.md` を厳密に踏襲する。
   - CRUD 画面なら、基本的に以下を揃える：
     - `GET /resource`
     - `POST /resource`
     - `GET /resource/{id}`
     - `PUT /resource/{id}`
     - 必要なら `DELETE /resource/{id}`

### 3-2. OpenAPI (`docs/openapi.yaml`) の自動生成

3. プロジェクトルートで、次のコマンドを**提案し、必要に応じて実行**する：

   ```bash
   make gen-openapi
   ```

   - これは `cmd/specgen` を使って `docs/specs/*.md` → `docs/openapi.yaml` を生成する。
   - 生成後、少なくとも **今回追加した画面に関する `paths` / `components.schemas`** を目視で確認する。

### 3-3. Goコード（型 & サーバーインターフェース）の自動生成

4. 同様に、次のコマンドを**提案し、必要に応じて実行**する：

   ```bash
   make gen-code
   ```

   - 内部的には `go generate ./...` を実行し、  
     `internal/api/openapi.gen.go` を更新する。
   - ここに、今回追加したエンドポイントの **request/response 型** と  
     **サーバーインターフェースのメソッド** が生えてくる。

### 3-4. Infrastructure/Persistence層の実装

5. `internal/infrastructure/persistence/{画面名}_repository.go` を作成 or 更新し、  
   SQLを実装する。

   - 例: `{画面名}Repository` struct を定義し、`pgxpool.Pool` を注入する。
   - CRUD であれば、`Create`, `FindByID`, `FindAll`, `FindByUserID`, `Update`, `Delete` のようなメソッドを実装する。
   - 共通パターンは `SubscriptionRepository` や `PlanRepository` を参考にする。

### 3-5. Usecase層の実装

6. `internal/usecase/{画面名}_service.go` を作成 or 更新し、  
   生成された型を使ってビジネスロジックを実装する。

   - 例: `{画面名}Service` struct を定義し、リポジトリを注入する。
   - CRUD であれば、`CreateX`, `GetX`, `ListX`, `UpdateX`, `DeleteX` のようなメソッドを持たせる。
   - 共通パターンは `SubscriptionService` や `PlanService` を参考にする。

### 3-6. Interface/HTTP/Handler層の実装

7. `internal/interface/http/handler/{画面名}_handler.go` を作成 or 更新し、  
   Echo ハンドラを実装する。

   - 生成されたサーバーインターフェースのメソッドと  
     1:1 で対応するハンドラを定義すること。
   - 認証が必要なエンドポイントには `internal/interface/http/middleware/auth.go` の JWT ミドルウェアを必ず通す。
   - `middleware.GetUserID` で `user_id` を取り出し、  
     **必ず user_id でフィルタされたCRUD** になるように usecase を呼び出す。

### 3-7. ルーティングの追加

8. `cmd/api/main.go` を更新し、`e.Group("/api")` の下に  
   新しい画面用のルートグループを追加する。

   - 例: `billing` 画面なら

     ```go
     billing := api.Group("/billing")
     billing.Use(authmw.AuthMiddleware(cfg.JWTSecret)) // 認証が必要な場合
     {
         billing.GET("", billingHandler.ListBilling)
         billing.POST("", billingHandler.CreateBilling)
         billing.GET("/:id", billingHandler.GetBilling)
         billing.PUT("/:id", billingHandler.UpdateBilling)
     }
     ```

   - ルーティングのパスは **OpenAPI (`docs/openapi.yaml`) の定義と完全一致**させること。

### 3-8. テストケースの生成

**🤖 AI向け:** 新しく追加したAPIのテストケースを生成する。  
`docs/API_TESTING.md` の形式に従って、以下の内容を含むテストケースファイルを作成する：

   - ファイル名: `docs/tests/{画面名}_test.sh` または `docs/tests/{画面名}_test.md`
   - 内容:
     - 認証が必要な場合: まずログインしてトークンを取得する手順
     - 各エンドポイント（GET, POST, PUT, DELETE）のcurlコマンド例
     - 成功時のレスポンス例
     - エラーケース（401, 400, 404など）のテスト
     - 必要な環境変数や前提条件

   **テストケースのテンプレート例:**

   ```bash
   #!/bin/bash
   # {画面名} API テストケース
   
   BASE_URL="http://localhost:8080/api"
   
   # 1. ログインしてトークンを取得
   TOKEN=$(curl -s -X POST $BASE_URL/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"testpass123"}' \
     | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
   
   echo "Token: ${TOKEN:0:50}..."
   
   # 2. GET /{resource} - 一覧取得
   echo "=== GET /{resource} ==="
   curl -s $BASE_URL/{resource} \
     -H "Authorization: Bearer $TOKEN" | jq .
   
   # 3. POST /{resource} - 作成
   echo "=== POST /{resource} ==="
   curl -s -X POST $BASE_URL/{resource} \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "field1": "value1",
       "field2": "value2"
     }' | jq .
   
   # 4. GET /{resource}/{id} - 取得
   echo "=== GET /{resource}/{id} ==="
   RESOURCE_ID="<作成時に取得したID>"
   curl -s $BASE_URL/{resource}/$RESOURCE_ID \
     -H "Authorization: Bearer $TOKEN" | jq .
   
   # 5. PUT /{resource}/{id} - 更新
   echo "=== PUT /{resource}/{id} ==="
   curl -s -X PUT $BASE_URL/{resource}/$RESOURCE_ID \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "field1": "updated_value1",
       "field2": "updated_value2"
     }' | jq .
   
   # 6. DELETE /{resource}/{id} - 削除（該当する場合）
   echo "=== DELETE /{resource}/{id} ==="
   curl -s -X DELETE $BASE_URL/{resource}/$RESOURCE_ID \
     -H "Authorization: Bearer $TOKEN"
   ```

   **注意事項:**
   - 認証が必要なエンドポイントには必ず `Authorization: Bearer $TOKEN` ヘッダーを含める
   - リクエストボディは `docs/specs/{画面名}.md` のデータモデル定義に基づいて作成する
   - エラーレスポンス（401, 400, 404, 500）のテストケースも含める
   - 実際のデータベースに影響を与えないよう、テスト用のデータを使用するか、テスト後にクリーンアップする

### 3-9. 最終確認

**🤖 AI向け:** 最低限、次を確認する：
   - `docs/specs/{画面名}.md` と `docs/openapi.yaml` に齟齬がないか。
   - 生成された `internal/api/openapi.gen.go` に  
     期待するメソッドと型が生えているか。
   - `cmd/api/main.go` のルーティングが、OpenAPIの `paths` と一致しているか。
   - テストケースが作成され、実際に動作確認できる状態になっているか。

---

## 4. この設計にするメリット

**🤖 AI向け:** 以下の点を意識して実装すること。

**👤 人間向け:** このセクションは、なぜこの設計にしたかの説明です。

- **追加が楽:**  
  新しい画面（例：請求履歴管理、タグ管理など）が欲しくなったら、  
  `docs/specs/*.md` に数ブロック追加して `make gen-all` するだけで、  
  APIと型の土台が一気に揃う。

- **仕様と実装のズレが減る:**  
  エンジニアもAIも **Markdown仕様書 → OpenAPI → コード生成** の一本線を守るため、  
  仕様書だけが古くなる、という状況を避けやすい。

- **AIとの相性が良い:**  
  ユーザーは「この画面を追加して」と伝えるだけで、  
  AIがこのファイルと `docs/WORKFLOW.md` を読み、  
  同じパターンで実装を増やしていける。

