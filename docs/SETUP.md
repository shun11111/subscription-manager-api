# 環境構築ガイド

## 📖 このドキュメントについて

このドキュメントでは、開発にジョインする人が必要なセットアップ手順を説明します。

**対象読者:**
- Go初心者の方
- 初めてこのプロジェクトに参加する方
- 環境構築を行う方

**読む前に:**
- ターミナル（コマンドライン）の基本的な操作ができること
- Gitがインストールされていること

**次のステップ:**
- 環境構築が完了したら → [`ENV_SETUP.md`](./ENV_SETUP.md) で環境変数を設定
- 環境変数設定が完了したら → [`ARCHITECTURE.md`](./ARCHITECTURE.md) でアーキテクチャを理解
- 開発を始める → [`WORKFLOW.md`](./WORKFLOW.md) で開発の流れを確認

---

このドキュメントでは、開発にジョインする人が必要なセットアップ手順を説明します。

## 前提条件

- ターミナル（コマンドライン）の基本的な操作ができること
- Gitがインストールされていること
- インターネット接続があること

## 1. Goのインストール

### macOS

#### 方法1: Homebrewを使う（推奨）

```bash
# Homebrewがインストールされていない場合
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Goをインストール
brew install go

# インストール確認
go version
```

#### 方法2: 公式インストーラーを使う

1. [Go公式サイト](https://go.dev/dl/) にアクセス
2. macOS用のインストーラー（`.pkg`ファイル）をダウンロード
3. ダウンロードしたファイルをダブルクリックしてインストール
4. ターミナルで確認:
   ```bash
   go version
   ```

### Windows

1. [Go公式サイト](https://go.dev/dl/) にアクセス
2. Windows用のインストーラー（`.msi`ファイル）をダウンロード
3. ダウンロードしたファイルをダブルクリックしてインストール
4. コマンドプロンプトまたはPowerShellで確認:
   ```bash
   go version
   ```

### Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# または公式のバイナリをインストール
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# PATHに追加（~/.bashrc または ~/.zshrc に追加）
export PATH=$PATH:/usr/local/go/bin

# インストール確認
go version
```

## 2. プロジェクトのクローン

```bash
# リポジトリをクローン
git clone <リポジトリのURL>
cd subscription-manager-api
```

## 3. 依存関係のインストール

```bash
# プロジェクトディレクトリで実行
go mod download
```

これで必要なパッケージが自動的にダウンロードされます。

## 4. 環境変数の設定

### 4-1. .envファイルの取得

**既存のメンバーから `.env` ファイルをもらってください。**

もらった `.env` ファイルをプロジェクトルート（`subscription-manager-api/` ディレクトリ）に配置します。

```bash
# .envファイルをプロジェクトルートに配置
# （既存メンバーからもらったファイルをそのまま使う）
```

**重要**: 
- `.env` ファイルには機密情報が含まれているため、**Gitにコミットしないでください**
- 既存メンバーからもらった `.env` ファイルをそのまま使用できます

## 5. サーバーの起動

### 5-1. 環境変数を読み込む

```bash
# .envファイルから環境変数を読み込む
export $(grep -v '^#' .env | xargs)
```

### 5-2. サーバーを起動

```bash
go run ./cmd/api/main.go
```

成功すると、以下のようなメッセージが表示されます：

```
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.12.0
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

これで `http://localhost:8080` でAPIが利用可能です。

## 6. 動作確認

別のターミナルを開いて、以下を実行：

```bash
# サインアップ
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

成功すると、`user_id` と `token` が返ってきます。

詳細なテスト方法は [`docs/API_TESTING.md`](API_TESTING.md) を参照してください。

## トラブルシューティング

### Goがインストールされていない

```bash
# エラーメッセージ: command not found: go
# → Goのインストールが必要です（上記の「1. Goのインストール」を参照）
```

### ポート8080が使用中

```bash
# エラーメッセージ: bind: address already in use
# 解決方法: 既存のプロセスを停止
lsof -ti:8080 | xargs kill -9

# Windowsの場合
netstat -ano | findstr :8080
taskkill /PID <PID番号> /F
```

### データベース接続エラー

```bash
# エラーメッセージ: failed to connect
# 確認事項:
# 1. .envファイルが正しく配置されているか（プロジェクトルート）
# 2. .envファイルのDATABASE_URLが正しいか
# 3. 環境変数が正しく読み込まれているか（echo $DATABASE_URL で確認）
```

### 環境変数が読み込まれない

```bash
# .envファイルを読み込む
export $(grep -v '^#' .env | xargs)

# 環境変数が設定されているか確認
echo $DATABASE_URL
echo $JWT_SECRET
```

### 依存関係のエラー

```bash
# 依存関係を再ダウンロード
go mod download

# 依存関係を整理
go mod tidy
```

## よくある質問

### Q: Goのバージョンは何が必要ですか？
A: Go 1.21以上が必要です。`go version` で確認できます。

### Q: .envファイルをGitにコミットしてもいいですか？
A: **いいえ、絶対にコミットしないでください**。機密情報が含まれているため、`.gitignore` で除外されています。


### Q: サーバーを停止する方法は？
A: ターミナルで `Ctrl + C` を押してください。

## 次のステップ

- [`docs/API_TESTING.md`](API_TESTING.md) - APIのテスト方法
- [`docs/ARCHITECTURE.md`](ARCHITECTURE.md) - アーキテクチャの説明
- [`docs/WORKFLOW.md`](WORKFLOW.md) - 開発ワークフロー

