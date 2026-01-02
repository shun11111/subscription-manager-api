# 統合開発フロー（フロントエンド + バックエンド）

## 📖 このドキュメントについて

このドキュメントは、**新画面を追加する際にAIが自動実行すべき完全なフロー**を定義しています。

**対象:**
- AI（Cursor、GitHub Copilot、Claude Codeなど）
- 新画面追加時の自動実装フロー

**重要:**
- このドキュメントは `subscription-manager-api` リポジトリにあります
- フロントエンド（`subscription-manager-admin`）の実装も含まれます
- 新画面追加時は、このドキュメントの手順を**必ず全て実行**してください

**関連ドキュメント:**
- [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) - API実装の詳細手順
- [`WORKFLOW.md`](./WORKFLOW.md) - API開発ワークフロー（人間向け）
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - アーキテクチャの説明

---

## 🎯 新画面追加時の完全フロー

ユーザーが「`{画面名}` という管理画面を追加したい」と指示した場合、以下の手順を**自動で実行**してください。

### 前提条件

1. **API側（subscription-manager-api）**
   - `docs/specs/{画面名}.md` が存在する、または作成される
   - 仕様書にはAPI仕様とフロント仕様の両方が含まれる

2. **フロント側（subscription-manager-admin）**
   - Next.js 16 (App Router) を使用
   - TypeScript + Tailwind CSS
   - 認証機能は `lib/auth.ts` を使用

---

## 📋 実装ステップ（全自動実行）

### フェーズ1: 仕様書の確認・作成

**🤖 AI向け:** まず、仕様書の存在を確認し、必要に応じて作成します。

1. `docs/specs/{画面名}.md` の存在を確認
2. 存在しない場合、`docs/specs/template.md` をベースに作成
3. 仕様書には以下を含める：
   - API仕様（エンドポイント、データモデル）
   - フロント仕様（画面構成、UI要件、ルーティング）

**仕様書テンプレート（拡張版）:**

```markdown
# {画面名}

## 概要
{画面の説明}

## API仕様

### エンドポイント
[既存のテンプレートに従って記述]

### データモデル
[既存のテンプレートに従って記述]

## フロント仕様

### 画面構成
- **ページパス**: `/{画面名}` または `/{画面名}/[id]`
- **レイアウト**: 一覧 / 詳細 / フォーム / 一覧+詳細
- **認証**: 必要 / 不要
- **ナビゲーション**: メニューに表示するか

### UI要件
- **一覧画面**（該当する場合）:
  - テーブル表示項目
  - ソート機能
  - フィルタ機能
  - ページネーション
  - 新規作成ボタン

- **詳細画面**（該当する場合）:
  - 表示項目
  - 編集ボタン
  - 削除ボタン

- **フォーム画面**（該当する場合）:
  - 入力項目とバリデーション
  - 必須/任意の区別
  - 送信ボタン

### ルーティング
- 一覧: `/{画面名}`
- 詳細: `/{画面名}/[id]`
- 新規作成: `/{画面名}/new`
- 編集: `/{画面名}/[id]/edit`
```

---

### フェーズ2: API実装（バックエンド）

**🤖 AI向け:** API側の実装を完了させます。詳細は [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) を参照してください。

#### 2-1. OpenAPI生成

```bash
cd subscription-manager-api
make gen-openapi
```

- `docs/openapi.yaml` が更新される
- 生成された内容を確認

#### 2-2. Goコード生成

```bash
make gen-code
```

- `internal/api/openapi.gen.go` が更新される
- 型とサーバーインターフェースが生成される

#### 2-3. Domain層の実装

- `internal/domain/{画面名}.go` を作成
- エンティティを定義

#### 2-4. Infrastructure/Persistence層の実装

- `internal/infrastructure/persistence/{画面名}_repository.go` を作成
- SQLクエリを実装
- 既存の `SubscriptionRepository` や `PlanRepository` を参考にする

#### 2-5. Usecase層の実装

- `internal/usecase/{画面名}_service.go` を作成
- ビジネスロジックを実装
- 既存の `SubscriptionService` や `PlanService` を参考にする

#### 2-6. Interface/HTTP/Handler層の実装

- `internal/interface/http/handler/{画面名}_handler.go` を作成
- Echoハンドラを実装
- 認証が必要な場合は `authmw.AuthMiddleware` を使用

#### 2-7. ルーティングの追加

- `cmd/api/main.go` を更新
- 新しいエンドポイントのルーティングを追加
- OpenAPIの定義と完全一致させる

#### 2-8. テストケースの生成

- `docs/tests/{画面名}_test.sh` または `docs/tests/{画面名}_test.md` を作成
- 各エンドポイントのcurlコマンド例を含める
- 認証が必要な場合はログイン手順を含める

---

### フェーズ3: フロントエンド実装

**🤖 AI向け:** フロントエンド側の実装を完了させます。`subscription-manager-admin` リポジトリで作業します。

#### 3-1. TypeScript型定義の生成（推奨）

OpenAPIからTypeScript型を自動生成する場合：

```bash
cd subscription-manager-admin

# openapi-typescript を使用する場合の例
npx openapi-typescript ../subscription-manager-api/docs/openapi.yaml -o src/types/api.ts
```

**注意:** 現在は手動で型定義を作成する場合もあります。

#### 3-2. APIクライアント関数の作成

- `lib/api/{画面名}.ts` を作成
- 各エンドポイントに対応する関数を実装
- 認証トークンを含める（`lib/auth.ts` の `getToken()` を使用）

**テンプレート例:**

```typescript
import { getToken } from '@/lib/auth';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export async function list{Resource}() {
  const token = getToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/{resource}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch');
  }

  return response.json();
}

export async function get{Resource}(id: string) {
  const token = getToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/{resource}/${id}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch');
  }

  return response.json();
}

export async function create{Resource}(data: Create{Resource}Request) {
  const token = getToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/{resource}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to create');
  }

  return response.json();
}

export async function update{Resource}(id: string, data: Update{Resource}Request) {
  const token = getToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/{resource}/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to update');
  }

  return response.json();
}

export async function delete{Resource}(id: string) {
  const token = getToken();
  if (!token) throw new Error('Not authenticated');

  const response = await fetch(`${API_URL}/{resource}/${id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to delete');
  }
}
```

#### 3-3. ページコンポーネントの作成

**一覧画面**（該当する場合）:
- `app/{画面名}/page.tsx` を作成
- データ取得、表示、新規作成ボタンを含める

**詳細画面**（該当する場合）:
- `app/{画面名}/[id]/page.tsx` を作成
- データ取得、表示、編集/削除ボタンを含める

**フォーム画面**（該当する場合）:
- `app/{画面名}/new/page.tsx` を作成（新規作成）
- `app/{画面名}/[id]/edit/page.tsx` を作成（編集）
- フォームバリデーションを含める

**テンプレート例（一覧画面）:**

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getToken } from '@/lib/auth';
import { list{Resource} } from '@/lib/api/{画面名}';

export default function {画面名}Page() {
  const router = useRouter();
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push('/login');
      return;
    }

    loadData();
  }, [router]);

  const loadData = async () => {
    try {
      setLoading(true);
      const result = await list{Resource}();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">{画面名}</h1>
        <button
          onClick={() => router.push('/{画面名}/new')}
          className="px-4 py-2 bg-blue-500 text-white rounded"
        >
          新規作成
        </button>
      </div>

      <table className="w-full border-collapse border border-gray-300">
        <thead>
          <tr>
            <th className="border border-gray-300 p-2">ID</th>
            {/* 他のカラム */}
          </tr>
        </thead>
        <tbody>
          {data.map((item) => (
            <tr key={item.id}>
              <td className="border border-gray-300 p-2">{item.id}</td>
              {/* 他のカラム */}
              <td className="border border-gray-300 p-2">
                <button
                  onClick={() => router.push(`/{画面名}/${item.id}`)}
                  className="text-blue-500"
                >
                  詳細
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

#### 3-4. UIコンポーネントの実装

- 必要に応じて `components/{画面名}/` ディレクトリを作成
- 再利用可能なコンポーネントを実装
- Tailwind CSSを使用してスタイリング

#### 3-5. 認証チェックの実装

- 認証が必要なページには、認証チェックを追加
- `lib/auth.ts` の `getToken()` を使用
- トークンがない場合は `/login` にリダイレクト

#### 3-6. エラーハンドリングの実装

- API呼び出し時のエラーハンドリング
- ユーザーフレンドリーなエラーメッセージの表示
- 401エラー時は自動的にログインページにリダイレクト

---

### フェーズ4: 統合テスト

**🤖 AI向け:** APIとフロントエンドの統合を確認します。

1. **APIサーバーの起動確認**
   ```bash
   cd subscription-manager-api
   go run ./cmd/api/main.go
   ```

2. **フロントエンドサーバーの起動確認**
   ```bash
   cd subscription-manager-admin
   npm run dev
   ```

3. **動作確認**
   - ブラウザで `http://localhost:3000` にアクセス
   - ログイン → 新画面にアクセス → CRUD操作をテスト

4. **エラーケースの確認**
   - 認証エラー（401）
   - バリデーションエラー（400）
   - 存在しないリソース（404）

---

### フェーズ5: 最終確認

**🤖 AI向け:** 以下の項目を確認してください。

#### API側の確認

- [ ] `docs/specs/{画面名}.md` が存在し、内容が正しい
- [ ] `docs/openapi.yaml` に新しいエンドポイントが追加されている
- [ ] `internal/api/openapi.gen.go` に型とインターフェースが生成されている
- [ ] Domain/Usecase/Infrastructure/Handler層が実装されている
- [ ] `cmd/api/main.go` にルーティングが追加されている
- [ ] テストケースが作成されている

#### フロント側の確認

- [ ] `lib/api/{画面名}.ts` が作成され、全エンドポイントに対応している
- [ ] ページコンポーネントが作成されている
- [ ] 認証チェックが実装されている
- [ ] エラーハンドリングが実装されている
- [ ] UIが仕様書の要件を満たしている

#### 統合確認

- [ ] APIとフロントエンドが正常に通信できる
- [ ] 認証フローが正しく動作する
- [ ] CRUD操作が全て動作する
- [ ] エラーケースが適切に処理される

---

## 🔄 既存画面を編集する場合

既存の画面を編集する場合も、同様のフローに従います：

1. `docs/specs/{画面名}.md` を編集
2. `make gen-openapi` でOpenAPIを更新
3. `make gen-code` でGoコードを更新
4. 必要に応じて各レイヤを修正
5. フロントエンド側も同様に更新

---

## 📝 注意事項

### API側

- 必ず `docs/specs/*.md` → `openapi.yaml` → `openapi.gen.go` の順で生成すること
- ルーティングのパスはOpenAPIの定義と完全一致させること
- 認証が必要なエンドポイントには必ず `authmw.AuthMiddleware` を使用すること
- `user_id` でフィルタされたCRUDになるように実装すること

### フロント側

- 認証が必要なページには必ず認証チェックを実装すること
- API呼び出し時は必ずエラーハンドリングを実装すること
- 環境変数 `NEXT_PUBLIC_API_URL` が設定されていることを確認すること
- 既存のUIパターン（`app/login/page.tsx` など）を参考にすること

### 統合

- APIサーバーとフロントエンドサーバーが両方起動していることを確認すること
- CORSの設定が必要な場合は確認すること
- 環境変数の設定を確認すること

---

## 🎯 まとめ

新画面追加時は、このドキュメントの手順を**必ず全て実行**してください。

1. 仕様書の確認・作成
2. API実装（フェーズ2）
3. フロントエンド実装（フェーズ3）
4. 統合テスト（フェーズ4）
5. 最終確認（フェーズ5）

各フェーズで生成・実装したファイルは、必ず確認してから次のフェーズに進んでください。

