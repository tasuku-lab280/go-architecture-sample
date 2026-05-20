# Layered Architecture Sample

DIP（依存性の逆転原則）を適用したレイヤードアーキテクチャの Go 実装サンプル。

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
│           ├── db.go            # *sql.DB の接続管理
│           └── user_repository.go  # user.Repository の MySQL 実装
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
