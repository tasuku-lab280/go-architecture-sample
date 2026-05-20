# Clean Architecture Sample

Uncle Bob のクリーンアーキテクチャを適用した Go 実装サンプル。

`layered/` との違いは、**ユースケースが入力ポート / 出力ポート（interface）を所有する**点。Controller は入力ポートに依存し、Presenter と Repository は出力ポートを実装する。

## ディレクトリ構成

```
clean/
├── cmd/api/                              # エントリポイント（DI 組み立て）
├── internal/
│   ├── entity/                           # 最内層: エンタープライズビジネスルール
│   │   └── user/                         #   User 集約（エンティティ + 値オブジェクト）
│   │       ├── user.go
│   │       ├── email.go                  #   Email 値オブジェクト
│   │       ├── password.go               #   Password 値オブジェクト
│   │       └── errors.go
│   │
│   ├── usecase/                          # アプリケーションビジネスルール
│   │   ├── port/
│   │   │   ├── input/                    #   入力ポート（Controller が依存）
│   │   │   │   └── register_user.go
│   │   │   └── output/                   #   出力ポート（Interactor が依存）
│   │   │       ├── user_repository.go    #     永続化境界
│   │   │       └── register_user_presenter.go  # 表示境界
│   │   └── interactor/                   #   ユースケース実装
│   │       └── register_user.go
│   │
│   ├── adapter/                          # インターフェースアダプタ層
│   │   ├── controller/                   #   HTTP → 入力ポート呼び出し
│   │   │   └── user_controller.go
│   │   ├── presenter/                    #   出力ポート実装、ViewModel 生成
│   │   │   └── register_user_presenter.go
│   │   ├── gateway/                      #   永続化ポート実装（MySQL）
│   │   │   └── user_repository.go
│   │   └── router/
│   │       └── router.go
│   │
│   └── infrastructure/                   # 最外層: フレームワーク・ドライバ
│       ├── database/                     #   *sql.DB の接続管理
│       │   └── mysql.go
│       └── server/                       #   net/http サーバ起動
│           └── http.go
│
├── db/init/                              # スキーマ定義 SQL
├── Dockerfile
└── compose.yaml
```

## 依存方向

```
infrastructure ──▶ adapter ──▶ usecase ──▶ entity
                                  ▲
                                  │
                       port/input ・ port/output
```

- `controller` は `usecase/port/input.RegisterUserInputPort` に依存（具象 Interactor は知らない）
- `gateway` は `usecase/port/output.UserRepository` を実装（DIP）
- `presenter` は `usecase/port/output.RegisterUserPresenter` を実装（DIP）
- `interactor` は出力ポート（interface）にだけ依存

## なぜ Presenter を分けるのか

`layered/` では Handler がユースケースの戻り値を直接 HTTP に書いていたが、Clean Architecture では Interactor が `Presenter.Present(output)` を呼び、Presenter 側で ViewModel（HTTP ステータスコードやレスポンスボディの形）を組み立てる。

これにより、**ドメイン例外と HTTP ステータスのマッピング**が Presenter 内に閉じ、ユースケースは「正常終了かエラーか」だけを知る純粋なオーケストレーションに留まる。

## リクエストごとに Presenter / Interactor を生成する理由

Presenter は ViewModel を内部に保持するためステートフル。並行リクエストで状態が混ざらないよう、Controller は **リクエストごとに** Presenter と Interactor を組み立てる。

ただし `Controller` が具象 Interactor を import すると依存方向が崩れる。そのため `cmd/api/main.go` でファクトリ関数 `func(output.RegisterUserPresenter) input.RegisterUserInputPort` を組み立てて Controller に注入し、Controller は「入力ポートを返すファクトリ」だけを知っている形にしている。

## 起動

```bash
docker compose up --build
curl -X POST http://localhost:8081/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@example.com","password":"password123"}'
```

`layered/` と同時起動するためポートは `8081` / DB は `3307` にずらしてある。

## テスト

```bash
# 単体テスト（entity / interactor / presenter / controller）
go test ./internal/entity/... ./internal/usecase/... \
        ./internal/adapter/controller/... ./internal/adapter/presenter/...

# 統合テスト（gateway、DB が必要）
docker compose up -d db
go test ./internal/adapter/gateway/...
```

### テスト戦略

| 層 | テスト種別 | テストダブル |
| --- | --- | --- |
| `entity/user` | 単体（純粋関数） | なし |
| `usecase/interactor` | 単体 | in-memory Repository / spy Presenter |
| `adapter/presenter` | 単体 | なし（出力データを直接渡す） |
| `adapter/controller` | 結合（HTTP〜Interactor） | in-memory Repository、Presenter は本物 |
| `adapter/gateway` | 統合（実 DB） | なし |
