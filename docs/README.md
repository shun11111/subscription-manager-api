# ドキュメント一覧

このディレクトリには、プロジェクトの開発・運用に必要なドキュメントが含まれています。

## 📚 初めての方へ

初めてこのプロジェクトに参加する方は、以下の順序でドキュメントを読むことをおすすめします：

### 1. まず読むべきドキュメント

1. **[SETUP.md](./SETUP.md)** - 環境構築ガイド
   - Goのインストール方法
   - プロジェクトのセットアップ手順
   - サーバーの起動方法
   - **対象**: Go初心者、初めてプロジェクトに参加する方

2. **[ENV_SETUP.md](./ENV_SETUP.md)** - 環境変数の設定
   - `.env`ファイルの取得方法
   - 必要な環境変数の説明
   - **対象**: 環境構築を行う方

### 2. プロジェクトの理解を深める

3. **[ARCHITECTURE.md](./ARCHITECTURE.md)** - アーキテクチャの説明
   - Clean Architectureのレイヤ構造
   - 各レイヤの役割と責務
   - リクエストの流れ
   - **対象**: コードを理解したい方、アーキテクチャを学びたい方

4. **[DIRECTORY_STRUCTURE.md](./DIRECTORY_STRUCTURE.md)** - ディレクトリ構成の説明
   - プロジェクトのディレクトリ構成
   - 各ディレクトリの役割
   - Clean Architectureとの関係
   - **対象**: プロジェクト構造を理解したい方

### 3. 開発を始める

5. **[WORKFLOW.md](./WORKFLOW.md)** - 開発ワークフロー
   - 新しい画面/APIを追加する手順
   - 仕様書の書き方
   - コード生成の流れ
   - **対象**: 新機能を追加する方

6. **[API_TESTING.md](./API_TESTING.md)** - APIのテスト方法
   - `curl`を使ったテスト
   - Postman/Insomniaの使い方
   - テストスクリプトの例
   - **対象**: APIをテストしたい方

## 📖 ドキュメント詳細

### セットアップ関連

- **[SETUP.md](./SETUP.md)**: 環境構築の完全ガイド（Go初心者向け）
- **[ENV_SETUP.md](./ENV_SETUP.md)**: 環境変数の設定方法

### アーキテクチャ・設計関連

- **[ARCHITECTURE.md](./ARCHITECTURE.md)**: Clean Architectureの説明と各レイヤの役割
- **[DIRECTORY_STRUCTURE.md](./DIRECTORY_STRUCTURE.md)**: ディレクトリ構成の詳細説明

### 開発関連

- **[WORKFLOW.md](./WORKFLOW.md)**: 新機能追加のワークフロー
- **[API_TESTING.md](./API_TESTING.md)**: APIのテスト方法

### 仕様書関連

- **[specs/](./specs/)**: Markdown形式のAPI仕様書
  - `auth.md`: 認証APIの仕様
  - `subscriptions.md`: サブスクリプションAPIの仕様
  - `plan_master.md`: プラン・価格マスターAPIの仕様
  - `template.md`: 仕様書のテンプレート

### API仕様

- **[openapi.yaml](./openapi.yaml)**: OpenAPI形式のAPI仕様書（自動生成）

## 🔍 よくある質問

### Q: どのドキュメントから読めばいいですか？

**A**: 初めての方は、以下の順序で読むことをおすすめします：
1. SETUP.md（環境構築）
2. ENV_SETUP.md（環境変数設定）
3. ARCHITECTURE.md（アーキテクチャ理解）
4. WORKFLOW.md（開発の流れ）

### Q: コードを理解したいのですが、どこから始めればいいですか？

**A**: 以下の順序で読むことをおすすめします：
1. ARCHITECTURE.md（全体像の理解）
2. DIRECTORY_STRUCTURE.md（ディレクトリ構成の理解）
3. 実際のコードを読む（`cmd/api/main.go`から始める）

### Q: 新しい機能を追加したいのですが、どうすればいいですか？

**A**: WORKFLOW.mdを参照してください。仕様書の作成からコード実装までの流れが説明されています。

### Q: APIをテストしたいのですが、どうすればいいですか？

**A**: API_TESTING.mdを参照してください。`curl`やPostmanを使ったテスト方法が説明されています。

## 📝 ドキュメントの更新

ドキュメントは、プロジェクトの変更に合わせて更新されます。疑問点や改善提案がある場合は、チームに共有してください。

## 🔗 関連リンク

- [プロジェクトルートのREADME.md](../README.md)
- [OpenAPI仕様書](./openapi.yaml)

