# JWT_SECRET の解説

## 📖 このドキュメントについて

このドキュメントでは、`JWT_SECRET`とは何か、なぜ必要なのか、どのように使われているのかを初心者向けに説明します。

**対象読者:**
- JWT_SECRETについて理解したい方
- 認証の仕組みを知りたい方
- セキュリティについて学びたい方

**関連ドキュメント:**
- [`ENV_SETUP.md`](./ENV_SETUP.md) - 環境変数の設定方法
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの説明

---

## JWT_SECRETとは？

`JWT_SECRET`は、**JWTトークンの署名と検証に使う秘密鍵**です。

### 簡単に言うと

```
JWT_SECRET = トークンの「印鑑」や「署名」を作るための秘密の鍵
```

---

## JWTトークンとは？

JWT（JSON Web Token）は、ユーザーがログインしたことを証明する「証明書」のようなものです。

### 例：ログインの流れ

```
1. ユーザーがログイン
   ↓
2. サーバーがJWTトークンを発行
   ↓
3. ユーザーがトークンを持ってAPIを呼び出す
   ↓
4. サーバーがトークンを検証して「この人は認証済み」と判断
```

### JWTトークンの構造

JWTトークンは3つの部分から構成されます：

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIiwgImV4cCI6MTYwOTQ1NjgwMH0.abc123def456...
│───────────────│ │──────────────────────────────│ │──────────────────────────│
   ヘッダー             ペイロード（データ）              署名（シグネチャ）
```

1. **ヘッダー**: トークンの種類や署名アルゴリズム
2. **ペイロード**: ユーザーIDや有効期限などの情報
3. **署名**: JWT_SECRETで作られた「印鑑」

---

## JWT_SECRETの役割

### 1. トークンの生成（署名）

ログイン時に、サーバーがJWTトークンを作成します：

```go
// internal/usecase/auth_service.go より
func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7日間有効
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))  // ← JWT_SECRETで署名
}
```

**イメージ：**
```
トークンの内容（ユーザーIDなど）
    ↓
JWT_SECRETで「印鑑」を押す
    ↓
完成したJWTトークン
```

### 2. トークンの検証（検証）

APIリクエスト時に、サーバーがトークンが本物か確認します：

```go
// internal/interface/http/middleware/auth.go より
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, jwt.ErrSignatureInvalid
    }
    return []byte(jwtSecret), nil  // ← JWT_SECRETで検証
})
```

**イメージ：**
```
受け取ったJWTトークン
    ↓
JWT_SECRETで「印鑑」を確認
    ↓
本物ならOK、偽物ならエラー
```

---

## なぜJWT_SECRETが必要なのか？

### セキュリティのため

JWT_SECRETがないと、誰でも偽のトークンを作れてしまいます。

#### ❌ JWT_SECRETがない場合

```
悪意のあるユーザー:
「自分はユーザーID=1だ」というトークンを勝手に作る
    ↓
サーバーがそれを信じてしまう（危険！）
```

#### ✅ JWT_SECRETがある場合

```
悪意のあるユーザー:
「自分はユーザーID=1だ」というトークンを作ろうとする
    ↓
でもJWT_SECRETがわからないので、正しい「印鑑」が押せない
    ↓
サーバーが「これは偽物だ」と判断して拒否（安全！）
```

### 比喩で理解する

JWT_SECRETは、**銀行の印鑑**のようなものです：

- **正しい印鑑（JWT_SECRET）**: サーバーだけが持っている
- **トークン**: 小切手のようなもの
- **検証**: 小切手に正しい印鑑が押されているか確認

---

## 実際のコードでの使用箇所

### 1. トークン生成（ログイン時）

```84:92:internal/usecase/auth_service.go
func (s *AuthService) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7日間有効
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
```

**何をしているか：**
- ユーザーIDと有効期限をトークンに含める
- JWT_SECRETで署名してトークンを完成させる

### 2. トークン検証（APIリクエスト時）

```35:40:internal/interface/http/middleware/auth.go
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
```

**何をしているか：**
- リクエストからJWTトークンを取得
- JWT_SECRETで署名を検証
- 本物ならリクエストを通す、偽物ならエラーを返す

---

## JWT_SECRETの設定方法

### 開発環境

`.env`ファイルに設定します：

```bash
JWT_SECRET=your-secret-key-change-this-in-production
```

### 本番環境

Fly.ioの場合：

```bash
flyctl secrets set JWT_SECRET="your-very-strong-secret-key"
```

### 安全なJWT_SECRETの生成方法

```bash
# ランダムな文字列を生成（推奨）
openssl rand -base64 32

# または
openssl rand -hex 32
```

**例：**
```bash
$ openssl rand -base64 32
aB3dEf9GhIjKlMnOpQrStUvWxYz1234567890+/=
```

---

## 重要な注意事項

### ⚠️ 絶対に守るべきこと

1. **Gitにコミットしない**
   - `.env`ファイルは`.gitignore`に含まれています
   - 絶対にGitHubにpushしないでください

2. **強力な文字列を使う**
   - 短すぎる文字列は危険です
   - 最低32文字以上を推奨

3. **本番環境では必ず変更**
   - 開発環境と本番環境で別の値を使う
   - デフォルト値は絶対に使わない

4. **定期的に変更**
   - 漏洩が疑われる場合は即座に変更
   - 変更すると、既存のトークンは全て無効になります

### ❌ やってはいけないこと

```bash
# ダメな例
JWT_SECRET=secret
JWT_SECRET=12345
JWT_SECRET=my-secret-key
JWT_SECRET=password
```

### ✅ 良い例

```bash
# 良い例
JWT_SECRET=aB3dEf9GhIjKlMnOpQrStUvWxYz1234567890+/=
JWT_SECRET=7f8a9b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2
```

---

## よくある質問

### Q: JWT_SECRETが漏洩したらどうなりますか？

**A**: 非常に危険です。攻撃者が正しいトークンを作成できるようになります。すぐに変更してください。

### Q: JWT_SECRETはどこに保存すればいいですか？

**A**: 
- **開発環境**: `.env`ファイル（Gitにコミットしない）
- **本番環境**: 環境変数やシークレット管理サービス（Fly.ioのsecrets、AWS Secrets Managerなど）

### Q: JWT_SECRETを変更するとどうなりますか？

**A**: 既存のトークンは全て無効になります。ユーザーは再度ログインする必要があります。

### Q: トークンに含まれる情報は暗号化されていますか？

**A**: いいえ、JWTトークンは**署名**されていますが、**暗号化**されていません。トークンの内容は誰でも見ることができます。機密情報（パスワードなど）は含めないでください。

### Q: なぜJWT_SECRETは環境変数で管理するのですか？

**A**: 
- コードに直接書くと、Gitにコミットされる危険がある
- 環境ごとに異なる値を使える
- コードを変更せずに設定を変更できる

---

## まとめ

| 項目 | 説明 |
|------|------|
| **JWT_SECRETとは** | JWTトークンの署名と検証に使う秘密鍵 |
| **役割** | トークンの生成（署名）と検証 |
| **なぜ必要** | 偽のトークンを作られないようにするため |
| **設定場所** | `.env`ファイル（開発）または環境変数（本番） |
| **注意点** | Gitにコミットしない、強力な文字列を使う |

---

## 次のステップ

- [`ENV_SETUP.md`](./ENV_SETUP.md) - 環境変数の設定方法
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの理解
- [`API_TESTING.md`](./API_TESTING.md) - APIのテスト方法

---

## 参考リンク

- [JWT公式サイト](https://jwt.io/)
- [JWT入門（日本語）](https://qiita.com/Naoto9282/items/8427918564400968bd68)

