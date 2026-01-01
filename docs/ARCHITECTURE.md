# アーキテクチャの説明

## 📖 このドキュメントについて

このドキュメントでは、プロジェクトのアーキテクチャと各レイヤの役割を説明します。

**対象読者:**
- コードを理解したい方
- アーキテクチャを学びたい方
- 新機能を追加する方

**読む前に:**
- 環境構築が完了していること（[`SETUP.md`](./SETUP.md) を参照）
- Goの基本的な構文を理解していること（推奨）

**関連ドキュメント:**
- [`DIRECTORY_STRUCTURE.md`](./DIRECTORY_STRUCTURE.md) - ディレクトリ構成の詳細
- [`WORKFLOW.md`](./WORKFLOW.md) - 開発ワークフロー

---

## アーキテクチャ概要

このプロジェクトは、シンプルな **Clean Architecture 風のレイヤ構造**を採用しています。

- **domain 層**: ビジネス上の「モノ」を表すエンティティ（構造体）とドメインルール
- **infrastructure/persistence 層**: DB との入出力（永続化）の実装
- **usecase 層**: ビジネスロジック（ユースケース）の実装
- **interface/http/handler 層**: HTTP リクエスト/レスポンス（Echo ハンドラ）

依存関係の向きは以下のようになります（内側ほどビジネス寄りで純粋）:

```text
interface → usecase → infrastructure → domain
```

`domain` はどこにも依存せず、もっとも内側のレイヤです。

---

## 各レイヤの役割

### `internal/domain`

**役割**: ビジネス上の概念（ユーザー、サブスクリプション、プランなど）を表す。

- 例: `User`, `Subscription`, `Plan` 構造体
- 共通フィールドは `BaseEntity` に定義（`ID`, `CreatedAt`, `UpdatedAt` など）
- HTTP や DB の具体的なこと（Echo の型、SQL 文字列など）はここには書かない

ポイント:

- 「現実世界の台帳」をコードで表現したもの、というイメージ
- 他のレイヤから再利用される「中心のモデル」

### `internal/infrastructure/persistence`

**役割**: Domain モデルを DB に保存・取得する。

- 例: `UserRepository`, `SubscriptionRepository`, `PlanRepository`
- やること:
  - SQL を組み立てて DB に対して `INSERT` / `SELECT` / `UPDATE` / `DELETE`
  - 取得した行を `domain.*` の構造体に読み込む (`Scan`)
- Usecase/Handler は SQL やテーブル名を直接知らず、Repository を経由して DB に触る

ポイント:

- 「`domain.Plan` ↔ `plans` テーブル 1行」の変換係
- DB を Postgres から他のものに変えたくなったとき、このレイヤを差し替えればよい設計

### `internal/usecase`

**役割**: ビジネスロジック（ユースケース）をまとめる。

- 例: `AuthService`, `SubscriptionService`, `PlanService`
- やること:
  - Repository を呼び出して、必要なデータを取得・保存
  - 必要に応じてバリデーションやドメインルールを実行
  - 複数の Repository を組み合わせた処理をここに書く
- HTTP/Echo の型は使わず、基本的に `context.Context` と Domain 型・Repository だけに依存する

ポイント:

- 「この操作を本当にやってよいか？どういう順番でやるか？」を決める頭脳部分
- Handler から見える「API としての入り口」に近いが、HTTP の詳細はここには持ち込まない

### `internal/interface/http/handler`

**役割**: HTTP の窓口（Echo ハンドラ）。  

- 例: `AuthHandler`, `SubscriptionHandler`, `PlanHandler`
- やること:
  1. `echo.Context` からリクエストを受け取る
  2. JSON を Go の構造体に `Bind` する
  3. Usecase を呼び出す
  4. Usecase の結果に応じて、HTTP ステータスコードと JSON を返す
- HTTP ステータスコード（200/201/400/401/404/500 など）とレスポンスフォーマットの責務を持つ

ポイント:

- 「URL + HTTP メソッド」ごとに**どの Usecase をどう呼ぶか**を定義する場所
- 認証ミドルウェアの適用や、ルーティングのグルーピング（`/api/subscriptions` など）と一緒に使われる

---

## 具体例: プラン作成フロー (`POST /api/plans`)

1. クライアントが `POST /api/plans` に JSON を送る
2. ルーティング（`cmd/api/main.go`）で `PlanHandler.CreatePlan` が呼ばれる
3. Handler:
   - JSON → `CreatePlanRequest` 構造体（handler 用 DTO）に `Bind`
   - `usecase.CreatePlanRequest` に詰め替えて `PlanService.CreatePlan` を呼ぶ
4. Usecase:
   - `CreatePlanRequest` から `domain.Plan` を構築
   - `PlanRepository.Create` を呼んで DB に保存
   - 保存した `*domain.Plan` を Handler に返す
5. Handler:
   - 201 (Created) として `domain.Plan` を JSON にシリアライズしてレスポンス

この流れの中で、

- **ビジネスの「プラン」という意味**は Domain/Usecase が担当
- **DB への INSERT/SELECT** は Infrastructure/Persistence が担当
- **HTTP リクエスト/レスポンス** は Interface/HTTP/Handler が担当  

というように責務を分担しています。

---

## 新しい画面/API を追加するときの基本手順

`docs/WORKFLOW.md` に詳しい手順がありますが、レイヤ構造に沿って見ると以下の流れになります。

1. **仕様を書く**  
   - `docs/specs/{画面名}.md` を作成（例: `plan_master.md`）
   - テンプレートに従ってエンドポイントとデータモデルを定義
2. **OpenAPI を更新**  
   - `make gen-openapi` で `docs/openapi.yaml` を生成/更新
3. **型・インターフェース生成**  
   - `make gen-code` で `internal/api/openapi.gen.go` を生成
4. **domain 層**  
   - `internal/domain/{resource}.go` にエンティティを追加
5. **infrastructure/persistence 層**  
   - `internal/infrastructure/persistence/{resource}_repository.go` を追加し、SQL を実装
6. **usecase 層**  
   - `internal/usecase/{resource}_service.go` を追加し、ビジネスロジックを実装
7. **interface/http/handler 層 & ルーティング**  
   - `internal/interface/http/handler/{resource}_handler.go` を追加
   - `cmd/api/main.go` にルーティングを追加

この流れに従うことで、新しい画面/API も既存のパターンに沿って実装できます。

---

## どこに何を書いてよいか迷ったときのガイド

- **「ビジネス用語が出てくるモデル」** → `internal/domain`
- **「SQL が出てくる」** → `internal/infrastructure/persistence`
- **「ユースケース名っぽいメソッド（CreateX, UpdateX, ListX）」** → `internal/usecase`
- **「HTTP ステータスコードや Echo の型が出てくる」** → `internal/interface/http/handler`

迷ったら、

- 「それは**誰にとっての都合の話**か？」（ユーザー/ビジネス？DB？HTTP？）

を考えて、その都合に一番近いレイヤに置く、というルールにするとチームで共有しやすくなります。


