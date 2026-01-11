---
trigger: always_on
glob: "**/*.go"
description: Go coding guidelines and best practices
---

# Go コーディングガイドライン

このドキュメントは、Sekishoプロジェクトのコーディング規約とベストプラクティスを概説します。[Effective Go](https://go.dev/doc/effective_go) の哲学と標準的なコミュニティの慣習に従います。

## 1. 全体的な哲学
- **シンプルさ**: コードはシンプルで理解しやすいものであるべきです。
- **可読性**: コードは書くよりも読まれることの方が多いです。読み手のために最適化してください。
- **保守性**: 明確なコードは保守しやすく、バグが発生しにくいです。

## 2. フォーマット
- **gofmt**: すべてのGoコードは `gofmt` を使用してフォーマットする必要があります。
- **Imports**: インポートの管理には `goimports` を使用してください。標準ライブラリとサードパーティライブラリのインポートを分けてグループ化します。

## 3. 命名規則
- **file_names.go**: ファイル名はスネークケースを使用します。
- **CamelCase**: 識別子（変数、関数、構造体）にはキャメルケースを使用します。
  - `ExportedIdentifier`（大文字で始まる）
  - `unexportedIdentifier`（小文字で始まる）
- **短い名前**: ローカル変数には短く簡潔な名前を使用します（例：`i`, `ctx`, `err`）。
- **説明的な名前**: エクスポートされる関数や型には、そのパッケージスコープに適した説明的な名前を使用します。
- **頭字語**: 頭字語は一貫性を持たせます（例：`ServeHttp` ではなく `ServeHTTP`）。

## 4. エラーハンドリング
- **明示的なチェック**: エラーは明示的に処理します。\ `_` を使用してエラーを無視しないでください。
- **ラッピング**: コンテキストを追加してエラーをラップする場合は、`fmt.Errorf("%w", err)` を使用します。
- **早期リターン**: エラーをチェックして早期にリターンすることで、正常系の処理のインデントを深くしないようにします。

```go
// Good
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
// proceed...

// Avoid
if err := doSomething(); err == nil {
    // proceed...
} else {
    return err
}
```

## 5. 並行処理
- **チャネル**: ゴルーチン間の連携やデータの受け渡しにはチャネルを使用します。
- **ミューテックス**: 共有状態の保護には `sync.Mutex` または `sync.RWMutex` を使用します。
- **Context**: キャンセル処理、タイムアウト、リクエストスコープの値の受け渡しには `context.Context` を使用します。関数の最初の引数として渡します。

```go
func DoSomething(ctx context.Context, arg string) error {
    // ...
}
```

## 6. プロジェクト構成
[Standard Go Project Layout](https://github.com/golang-standards/project-layout) に従います：

- `cmd/`: メインアプリケーション（エントリーポイント）。
- `internal/`: プライベートなアプリケーションおよびライブラリコード。
- `pkg/`: 他のプロジェクトからインポートしても安全なライブラリコード。
- `api/`: OpenAPI/Swagger仕様、JSONスキーマファイル、プロトコル定義ファイル。

## 7. 設定
- 設定には環境変数を使用します（`kelseyhightower/envconfig` や `joho/godotenv` などのライブラリ、または単純な場合は標準ライブラリを使用）。
- シークレットや環境固有の値をハードコードしないでください。

## 8. テスト
- **テーブル駆動テスト**: 複数のシナリオを簡潔に網羅するためにテーブル駆動テストを使用します。
- **並列テスト**: 意味のあるテストでは `t.Parallel()` を呼び出すべきです。
- **テストパッケージ**: テストを分離するために適切な場合は、外部テスト（ブラックボックステスト）に `_test` パッケージを使用します。

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -1, -2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 9. Linting
- Lintingには `golangci-lint` を使用します。
- 有効にすべき一般的なLinter: `govet`, `staticcheck`, `errcheck`, `gocritic`.

## 10. ドキュメント
- **コメント**: エクスポートされた関数や型にはコメントが必要です。
- **Godoc**: Godocと互換性のある形式でコメントを記述します。コメントは対象の名前で始めてください。

```go
// User represents a user in the system.
type User struct { ... }
```
