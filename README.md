# kanboard

Go言語用の [Kanboard](https://kanboard.org/) JSON-RPC APIクライアントライブラリです。

## 概要

Kanboard の JSON-RPC API を簡潔に呼び出すためのリクエスト生成・送信ユーティリティを提供します。ジェネリクスを活用した型安全なリクエスト/レスポンス処理が特徴です。

## 必要環境

- Go 1.18 以上（ジェネリクス使用）
- Kanboard サーバー（セルフホスト）

## インストール

```bash
go get github.com/sg-daigo/kanboard-go
```

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `KB_TOKEN` | Kanboard API トークン |

## 使い方

### サーバー設定

```go
// デフォルト設定（localhost + 環境変数からトークンを読み込み）
conf := kanboard.NewServerConfig()

// カスタム設定
conf := kanboard.ServerConfig{
    Server: "https://kanboard.example.com",
    Token:  "your-api-token",
}
```

### タスクのタグを取得する

```go
req := kanboard.NewTaskTagsRequest(42) // タスクID

result, err := kanboard.SendRequest[kanboard.TaskTagsParams, kanboard.TaskTagsResult](
    req, conf, nil,
)
if err != nil {
    log.Fatal(err)
}

for id, name := range result {
    fmt.Printf("Tag %s: %s\n", id, name)
}
```

### 任意のAPIメソッドを呼び出す

`NewRequest` と `SendRequest` を組み合わせることで、任意の Kanboard API メソッドを呼び出せます。

```go
type GetProjectParams struct {
    ProjectID int `json:"project_id"`
}

req := kanboard.NewRequest[GetProjectParams]("getProjectById", &GetProjectParams{
    ProjectID: 1,
})

project, err := kanboard.SendRequest[GetProjectParams, kanboard.ProjectResult](
    req, conf, nil,
)
if err != nil {
    log.Fatal(err)
}

fmt.Println(project.Name)
```

### ロギング

`SendRequest` / `SendRequestRaw` の第3引数に `*slog.Logger` を渡すと、レスポンスボディを DEBUG レベルで出力します。

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

result, err := kanboard.SendRequest[kanboard.TaskTagsParams, kanboard.TaskTagsResult](
    req, conf, logger,
)
```

## API リファレンス

### 型

| 型 | 説明 |
|----|------|
| `ServerConfig` | サーバーURLとAPIトークンの設定 |
| `Request[T]` | JSON-RPC リクエスト（ジェネリクス） |
| `Response` | JSON-RPC レスポンス |
| `FlexibleInt` | 数値・文字列どちらの形式でもデシリアライズ可能な整数型 |
| `TaskTagsParams` | `getTaskTags` メソッドのパラメータ |
| `TaskTagsResult` | `getTaskTags` メソッドの結果（`map[string]string`） |
| `ProjectResult` | プロジェクト情報 |
| `AllTagsResponse` | タグ情報 |

### 関数

| 関数 | 説明 |
|------|------|
| `NewServerConfig()` | デフォルトのサーバー設定を生成 |
| `NewRequest[T](method, params)` | 汎用 JSON-RPC リクエストを生成 |
| `NewTaskTagsRequest(taskID)` | `getTaskTags` リクエストを生成 |
| `SendRequestRaw[T](req, conf, logger)` | リクエストを送信し、生の `Response` を返す |
| `SendRequest[P, R](req, conf, logger)` | リクエストを送信し、結果を型 `R` にデシリアライズして返す |

## ライセンス

MIT
