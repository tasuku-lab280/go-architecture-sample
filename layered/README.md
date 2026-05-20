# Layered Architecture Sample

DIP（依存性の逆転原則）を適用したレイヤードアーキテクチャの Go 実装サンプル。ユーザー登録機能を題材に、各レイヤーの責務分離・テスト戦略・値オブジェクト設計を実演する。

## 技術スタック

- Go 1.25
- MySQL 8.4
- `net/http`（標準ライブラリ）
- bcrypt（パスワードハッシュ）

## ディレクトリ構成

```
layered/
├── cmd/api/                     # エントリポイント（DI 組み立て）
├── internal/
│   ├── handler/                 # HTTP リクエスト・レスポンス処理
│   │   ├── router.go
│   │   └── user.go
│   ├── usecase/                 # アプリケーションロジック（オーケストレーション）
│   │   └── register_user.go
│   ├── domain/
│   │   └── user/                # User 集約
│   │       ├── email.go         # Email 値オブジェクト
│   │       ├── password.go      # Password 値オブジェクト
│   │       ├── user.go          # User エンティティ
│   │       ├── repository.go    # Repository interface（DIP）
│   │       └── errors.go        # ドメインエラー
│   └── infrastructure/
│       └── database/            # MySQL Repository 実装
├── db/init/                     # スキーマ定義 SQL
├── Dockerfile
└── compose.yaml
```

## 依存方向

```
handler ──▶ usecase ──▶ domain ◀── infrastructure
                          ▲
                          │
                    Repository interface
```

`infrastructure → domain` の方向に矢印が逆転しているのが DIP のポイント。ドメインは技術詳細（DB ドライバ等）を一切知らない。

## 起動

```bash
docker compose up
```

- API: `http://localhost:8080`
- MySQL: `localhost:3306`

## API

### `POST /users`

ユーザー登録。

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

| Status | 条件 |
|---|---|
| 201 | 登録成功 |
| 400 | バリデーションエラー（email 形式不正・password 短すぎ） |
| 409 | email 重複 |
| 500 | サーバーエラー |

## テスト

```bash
# 全テスト
docker compose run --rm app go test ./...

# レイヤー別
docker compose run --rm app go test ./internal/domain/user/...
docker compose run --rm app go test ./internal/usecase/...
docker compose run --rm app go test ./internal/handler/...
docker compose run --rm app go test ./internal/infrastructure/database/...
```

### テスト戦略

| レイヤー | テスト方式 | 備考 |
|---|---|---|
| `domain` | 単体テスト | 外部依存なし、高速 |
| `usecase` | Fake repository を注入 | DB 不要、振る舞いベース検証 |
| `handler` | Fake repository を経由した E2E 風 | `httptest` で HTTP レイヤーごと検証 |
| `infrastructure` | 実 MySQL に対する統合テスト | docker compose 必須 |

## 設計判断

- **DIP**: Repository interface は domain 層に配置し、infrastructure が実装する。これによりドメインを技術詳細から独立させる
- **値オブジェクト**: `Email` / `Password` はプリミティブ型でなく専用型として定義し、バリデーションと振る舞いを型に集約する
- **ドメインエラー**: 業務ルール違反（email 重複・形式不正・password 短すぎ）は domain パッケージで定義し、handler 層で HTTP ステータスにマッピングする
- **集約サブパッケージ化**: `domain/` 配下を集約単位（`user/`）で分割し、将来 `order/` などを追加してもスケールする構造にする
- **テストダブルは Fake 中心**: gomock 等の生成モックではなく、手書きの in-memory 実装で振る舞いベースのテストを行う
