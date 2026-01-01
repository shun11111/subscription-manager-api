# 環境変数の設定

## 📖 このドキュメントについて

このドキュメントでは、環境変数の設定方法を説明します。

**対象読者:**
- 環境構築を行う方
- `.env`ファイルを設定する方

**読む前に:**
- Goのインストールが完了していること（[`SETUP.md`](./SETUP.md) を参照）

**次のステップ:**
- 環境変数設定が完了したら → [`ARCHITECTURE.md`](./ARCHITECTURE.md) でアーキテクチャを理解
- サーバーを起動してテスト → [`API_TESTING.md`](./API_TESTING.md) でAPIをテスト

---方法

## 概要

このプロジェクトでは、環境変数は**個人ごとに設定**します。機密情報を含むため、`.env`ファイルはGitにコミットしません。

## セットアップ手順

### 1. `.env.example`をコピー

```bash
cp .env.example .env
```

### 2. `.env`ファイルを編集

各環境変数に実際の値を設定してください（1Passwordなどから取得）。

```bash
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgres://postgres:your_password@your_project_ref.supabase.co:5432/postgres?sslmode=require
JWT_SECRET=your-secret-key-change-this-in-production
```

### 3. 環境変数を読み込む

`.env`ファイルを読み込むには、以下のいずれかの方法を使用してください：

#### 方法1: exportコマンドで手動設定

```bash
export PORT=8080
export DATABASE_URL="postgres://..."
export JWT_SECRET="your-secret"
```

#### 方法2: .envファイルを読み込む（シェルスクリプト）

`.env`ファイルに`export`文を書いて、`source`で読み込む：

```bash
# .envファイルの内容
export PORT=8080
export DATABASE_URL="postgres://..."
export JWT_SECRET="your-secret"

# 読み込み
source .env
```

## 環境変数の取得方法

### Supabase の接続情報

1. [Supabase Dashboard](https://app.supabase.com/) にログイン
2. プロジェクトを選択
3. **Settings** > **Database** に移動
4. **Connection string** セクションから接続文字列をコピー
   - **Connection pooling** ではなく、**Direct connection** を使用
   - パスワードは **Database password** を参照

### JWT_SECRET

JWT_SECRETは、JWTトークンの署名と検証に使う秘密鍵です。

**詳細な説明は [`JWT_SECRET.md`](./JWT_SECRET.md) を参照してください。**

開発環境では任意の文字列でOKです。本番環境では強力なランダム文字列を使用してください。

```bash
# ランダムなシークレットを生成（推奨）
openssl rand -base64 32

# または
openssl rand -hex 32
```

**重要**: 
- Gitにコミットしないでください
- 最低32文字以上の強力な文字列を使用してください
- 本番環境では必ず開発環境とは別の値を使用してください

## 環境変数の共有方法

### 推奨：チーム共有ドキュメント

機密情報は以下のいずれかで共有してください：

- **1Password / Bitwarden** などのパスワード管理ツール
- **Notion / Google Docs** などの共有ドキュメント（アクセス制限あり）
- **Slack のプライベートチャンネル**（一時的な共有のみ）

⚠️ **注意**: `.env`ファイルや環境変数の値を**GitHubのIssueやPRに直接貼り付けない**でください。

### 開発環境ごとの設定

- **個人開発環境**: 各自のSupabaseプロジェクト（またはローカルPostgreSQL）に接続
- **ステージング環境**: チーム共有のSupabaseプロジェクトに接続（接続情報は上記の方法で共有）
- **本番環境**: CI/CDのシークレット管理機能を使用

## トラブルシューティング

### 環境変数が読み込まれない

- `echo $DATABASE_URL`で環境変数が設定されているか確認
- `.env`ファイルが正しい場所にあるか確認（プロジェクトルート）
- 新しいターミナルでは設定が消えるので、再度`source .env`または`export`が必要

### Supabase接続エラー

- SSLモードを確認（Supabaseは`sslmode=require`が必要）
- パスワードが正しいか確認
- プロジェクトのIPアドレスが許可されているか確認（Supabase Dashboard > Settings > Database > Connection pooling）

