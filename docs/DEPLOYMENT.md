# デプロイメントガイド

## 📖 このドキュメントについて

このドキュメントでは、Dockerを使ったアプリケーションのビルドとデプロイ方法を説明します。

**対象読者:**
- アプリケーションを本番環境にデプロイする方
- Dockerの使い方を学びたい方
- Fly.ioへのデプロイを行う方

**読む前に:**
- Goの基本的な知識があること
- Dockerの基本的な概念を理解していること（このドキュメントで説明します）

**次のステップ:**
- デプロイが完了したら → [`API_TESTING.md`](./API_TESTING.md) でAPIをテスト

---

## Dockerとは？

Dockerは、アプリケーションを「コンテナ」という箱に入れて、どこでも同じ環境で動かせるようにするツールです。

### なぜDockerを使うのか？

#### 問題：環境の違いによる不具合

開発環境と本番環境で環境が違うと、不具合が起きやすいです。

```
開発環境（Mac）         本番環境（Linuxサーバー）
Go 1.23                Go 1.21（古い）
動く ✅                動かない ❌
```

#### 解決策：Dockerを使う

Dockerは「同じ環境を再現する箱」を作ります。

```
開発環境（Mac）         本番環境（Linuxサーバー）
┌─────────────┐        ┌─────────────┐
│   Docker    │        │   Docker    │
│  ┌───────┐  │        │  ┌───────┐  │
│  │ Go    │  │        │  │ Go    │  │
│  │ 1.23  │  │  同じ  │  │ 1.23  │  │
│  │       │  │  ←→   │  │       │  │
│  │ アプリ│  │        │  │ アプリ│  │
│  └───────┘  │        │  └───────┘  │
└─────────────┘        └─────────────┘
```

### Dockerfileとは？

Dockerfileは、その「箱」の作り方を書いた設計図です。

---

## Dockerfileの解説

このプロジェクトのDockerfileは以下の通りです：

```dockerfile
# 1. Goのバージョンを指定（1.23を使う）
ARG GO_VERSION=1.23

# 2. ビルド用の箱を作る（Goが入っている）
FROM golang:${GO_VERSION}-bookworm as builder

# 3. 作業ディレクトリを設定
WORKDIR /usr/src/app

# 4. 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 5. 依存関係をダウンロード
RUN go mod download && go mod verify

# 6. ソースコードを全部コピー
COPY . .

# 7. アプリケーションをビルド（実行ファイルを作る）
RUN go build -v -o /run-app ./cmd/api

# 8. 実行用の箱を作る（軽量なLinux）
FROM debian:bookworm

# 9. ビルドした実行ファイルだけをコピー
COPY --from=builder /run-app /usr/local/bin/run-app

# 10. 起動コマンドを指定
CMD ["/usr/local/bin/run-app"]
```

### 各行の説明

1. **ARG GO_VERSION=1.23**: Goのバージョンを指定（変更可能）
2. **FROM golang:...**: Goが入ったベースイメージを使用
3. **WORKDIR**: 作業ディレクトリを設定
4. **COPY go.mod go.sum**: 依存関係ファイルをコピー（キャッシュ効率化）
5. **RUN go mod download**: 依存関係をダウンロード
6. **COPY . .**: ソースコードをコピー
7. **RUN go build**: アプリケーションをビルド
8. **FROM debian:bookworm**: 軽量な実行環境を作成
9. **COPY --from=builder**: ビルドした実行ファイルだけをコピー
10. **CMD**: 起動時に実行するコマンド

---

## ローカルでDockerを使う

### 1. Dockerイメージをビルド

```bash
# プロジェクトルートで実行
docker build -t subscription-manager-api .
```

- `-t subscription-manager-api`: イメージに名前を付ける
- `.`: カレントディレクトリをビルドコンテキストにする

### 2. コンテナを起動

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname" \
  -e JWT_SECRET="your-secret-key" \
  -e PORT="8080" \
  subscription-manager-api
```

- `-p 8080:8080`: ホストの8080ポートをコンテナの8080ポートにマッピング
- `-e`: 環境変数を設定

### 3. バックグラウンドで実行

```bash
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="your-secret" \
  --name subscription-api \
  subscription-manager-api
```

- `-d`: バックグラウンドで実行
- `--name`: コンテナに名前を付ける

### 4. コンテナの確認・停止

```bash
# 実行中のコンテナを確認
docker ps

# コンテナを停止
docker stop subscription-api

# コンテナを削除
docker rm subscription-api
```

---

## Fly.ioへのデプロイ

Fly.ioは、Dockerコンテナを簡単にデプロイできるプラットフォームです。

### 前提条件

1. Fly.ioアカウントを作成: https://fly.io
2. Fly CLIをインストール: https://fly.io/docs/hands-on/install-flyctl/

```bash
# macOS
brew install flyctl

# または公式インストーラー
curl -L https://fly.io/install.sh | sh
```

### デプロイ手順

#### 1. Fly.ioにログイン

```bash
flyctl auth login
```

#### 2. アプリを作成（初回のみ）

```bash
flyctl launch
```

このコマンドで以下が作成されます：
- `fly.toml`: Fly.ioの設定ファイル
- アプリケーションの登録

#### 3. 環境変数を設定

```bash
# データベースURL
flyctl secrets set DATABASE_URL="postgres://..."

# JWTシークレット
flyctl secrets set JWT_SECRET="your-secret-key"

# ポート（通常は8080）
flyctl secrets set PORT="8080"
```

または、Fly.ioのWebダッシュボードから設定することもできます。

#### 4. デプロイ

```bash
# GitHubと連携している場合、自動的にデプロイされます
# 手動でデプロイする場合
flyctl deploy
```

このコマンドで以下が実行されます：
1. GitHubからコードを取得
2. Dockerfileを読み込む
3. Dockerイメージをビルド
4. コンテナを起動
5. インターネットに公開

#### 5. デプロイの確認

```bash
# アプリの状態を確認
flyctl status

# ログを確認
flyctl logs

# アプリのURLを確認
flyctl open
```

---

## .dockerignoreファイル

`.dockerignore`は、Dockerビルド時に除外するファイルを指定します。

このプロジェクトの`.dockerignore`は以下の通りです：

```
# Exclude generate.go from Docker build
generate.go

# Exclude other build artifacts
*.exe
*.test
*.out
main

# Exclude environment files
.env
.env.local

# Exclude IDE files
.idea/
.vscode/
*.swp
*.swo

# Exclude git
.git/
.gitignore

# Exclude documentation (optional, but reduces build context)
docs/

# Exclude test files (optional)
*_test.go
```

### なぜ必要？

- **ビルド速度の向上**: 不要なファイルをコピーしない
- **セキュリティ**: `.env`などの機密情報を誤って含めない
- **イメージサイズの削減**: 必要なファイルだけを含める

---

## トラブルシューティング

### ビルドエラー: "os imported and not used"

**原因**: `generate.go`がビルドに含まれている

**解決方法**: `.dockerignore`に`generate.go`を追加（既に追加済み）

### デプロイエラー: "failed to connect to database"

**原因**: 環境変数が設定されていない

**解決方法**: 
```bash
flyctl secrets set DATABASE_URL="postgres://..."
```

### ポートエラー: "port already in use"

**原因**: ローカルで既にポート8080が使用されている

**解決方法**: 
```bash
# 使用中のプロセスを確認
lsof -ti:8080

# プロセスを停止
lsof -ti:8080 | xargs kill -9
```

### Dockerイメージが大きい

**原因**: 不要なファイルが含まれている

**解決方法**: 
- `.dockerignore`を確認
- マルチステージビルドを使用（既に使用中）

---

## よくある質問

### Q: ローカルで開発する時もDockerを使うべきですか？

**A**: 開発時は`go run ./cmd/api/main.go`を使うことをおすすめします。Dockerは主にデプロイ時に使用します。

### Q: Dockerfileを変更したらどうすればいいですか？

**A**: 
1. 変更をGitHubにpush
2. `flyctl deploy`を実行（または自動デプロイを待つ）

### Q: 環境変数はどこで設定すればいいですか？

**A**: 
- **ローカル開発**: `.env`ファイル
- **Fly.io**: `flyctl secrets set`コマンドまたはWebダッシュボード

### Q: 複数の環境（開発・ステージング・本番）を管理するには？

**A**: Fly.ioでは、別々のアプリを作成するか、同じアプリで環境変数を切り替えます。

---

## 次のステップ

- [`API_TESTING.md`](./API_TESTING.md) - デプロイ後のAPIテスト
- [`WORKFLOW.md`](./WORKFLOW.md) - 開発ワークフロー
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの理解

---

## 参考リンク

- [Docker公式ドキュメント](https://docs.docker.com/)
- [Fly.io公式ドキュメント](https://fly.io/docs/)
- [Go公式ドキュメント](https://go.dev/doc/)

