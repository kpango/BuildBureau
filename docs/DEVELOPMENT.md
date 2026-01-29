# Development Guide

BuildBureauの開発者向けガイド

## 開発環境のセットアップ

### 必要なツール

1. **Go 1.23以上**
   ```bash
   go version
   ```

2. **Protocol Buffers compiler (protoc)**
   ```bash
   # macOS
   brew install protobuf
   
   # Linux
   sudo apt install -y protobuf-compiler
   ```

3. **Go用のprotocプラグイン**
   ```bash
   make install-tools
   ```

### プロジェクトのクローン

```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
```

### 依存関係のインストール

```bash
make deps
```

## ビルド

### 通常ビルド

```bash
make build
```

ビルドされたバイナリは`./bin/buildbureau`に配置されます。

### クリーンビルド

```bash
make clean
make build
```

## テスト

### 全テスト実行

```bash
make test
```

### 特定パッケージのテスト

```bash
go test -v ./internal/config
go test -v ./internal/agent
```

### テストカバレッジ

```bash
go test -cover ./...
```

## コード品質

### フォーマット

```bash
make fmt
```

### Lint

```bash
make vet
```

または統合コマンド：

```bash
make lint
```

## プロジェクト構造

```
BuildBureau/
├── cmd/
│   └── buildbureau/          # メインアプリケーション
│       └── main.go
├── internal/                  # 内部パッケージ
│   ├── agent/                # エージェント実装
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── config/               # 設定管理
│   │   ├── config.go
│   │   └── config_test.go
│   ├── grpc/                 # gRPCサービス実装
│   ├── slack/                # Slack通知
│   │   └── notifier.go
│   └── ui/                   # Terminal UI
│       └── ui.go
├── proto/                     # Protocol Buffers定義
│   └── buildbureau/v1/
│       └── service.proto
├── pkg/                       # 公開パッケージ
│   ├── models/               # データモデル
│   └── utils/                # ユーティリティ
├── docs/                      # ドキュメント
│   ├── ARCHITECTURE.md
│   └── CONFIGURATION.md
├── config.yaml               # デフォルト設定
├── .env.example              # 環境変数テンプレート
├── Makefile                  # ビルドスクリプト
├── go.mod                    # Go依存関係
└── README.md
```

## 新機能の追加

### 1. 新しいエージェント実装

```go
// internal/agent/president.go
package agent

import (
    "context"
    "github.com/kpango/BuildBureau/internal/config"
)

type PresidentAgent struct {
    *BaseAgent
    // 追加のフィールド
}

func NewPresidentAgent(id string, cfg config.AgentConfig) *PresidentAgent {
    return &PresidentAgent{
        BaseAgent: NewBaseAgent(id, AgentTypePresident, cfg),
    }
}

func (a *PresidentAgent) Process(ctx context.Context, input interface{}) (interface{}, error) {
    // 実装
    return nil, nil
}
```

### 2. 新しいgRPCサービスの追加

1. `proto/buildbureau/v1/service.proto`にサービス定義を追加

```protobuf
service NewService {
    rpc NewMethod(RequestType) returns (ResponseType);
}
```

2. プロトコルバッファのコード生成

```bash
make proto
```

3. サービス実装を追加

```go
// internal/grpc/new_service.go
package grpc

type NewServiceServer struct {
    // フィールド
}

func (s *NewServiceServer) NewMethod(ctx context.Context, req *pb.RequestType) (*pb.ResponseType, error) {
    // 実装
    return &pb.ResponseType{}, nil
}
```

### 3. 新しいSlack通知イベントの追加

1. `internal/config/config.go`に通知設定を追加

```go
type NotificationsConfig struct {
    // 既存のイベント
    NewEvent NotificationConfig `yaml:"newEvent"`
}
```

2. `config.yaml`に設定を追加

```yaml
slack:
  notifications:
    newEvent:
      enabled: true
      message: "新しいイベント: {{.Data}}"
```

3. `internal/slack/notifier.go`にメソッドを追加

```go
func (n *Notifier) SendNewEvent(ctx context.Context, data string) error {
    return n.Send(ctx, NotificationNewEvent, NotificationData{
        // データ
    })
}
```

## デバッグ

### ログレベルの設定

```yaml
# config.yaml
ui:
  logLevel: "debug"
```

または環境変数：

```bash
export LOG_LEVEL=debug
```

### デバッガの使用

```bash
# Delveのインストール
go install github.com/go-delve/delve/cmd/dlv@latest

# デバッグ実行
dlv debug ./cmd/buildbureau
```

## テストの書き方

### ユニットテスト

```go
// internal/agent/president_test.go
package agent

import (
    "context"
    "testing"
)

func TestPresidentAgent_Process(t *testing.T) {
    agent := NewPresidentAgent("test-1", config.AgentConfig{})
    
    ctx := context.Background()
    result, err := agent.Process(ctx, "test input")
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    // アサーション
}
```

### テーブル駆動テスト

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"case1", "input1", "output1", false},
        {"case2", "input2", "output2", false},
        {"error case", "bad", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := functionUnderTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## プルリクエストのガイドライン

### ブランチ戦略

- `main`: 本番ブランチ
- `develop`: 開発ブランチ
- `feature/*`: 機能追加ブランチ
- `fix/*`: バグ修正ブランチ
- `docs/*`: ドキュメント更新ブランチ

### コミットメッセージ

Conventional Commitsに従う：

```
<type>(<scope>): <subject>

<body>

<footer>
```

タイプ：
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: コードスタイル
- `refactor`: リファクタリング
- `test`: テスト
- `chore`: ビルド・設定等

例：
```
feat(agent): Add LLM integration for president agent

Implement LLM API calls using Google ADK.
Add retry logic and error handling.

Closes #123
```

### プルリクエストチェックリスト

- [ ] コードがビルドできる
- [ ] 全テストが通る
- [ ] 新しいテストを追加した
- [ ] ドキュメントを更新した
- [ ] コミットメッセージが適切
- [ ] コードレビューを受けた

## よくある問題と解決策

### 1. ビルドエラー: "cannot find package"

```bash
make deps
go mod tidy
```

### 2. テストエラー: "no such file or directory"

パスが相対パスになっていないか確認：

```go
// 悪い例
os.ReadFile("config.yaml")

// 良い例
os.ReadFile("/absolute/path/config.yaml")
```

### 3. Protocol Buffer生成エラー

```bash
make install-tools
make proto
```

## パフォーマンスプロファイリング

### CPUプロファイル

```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### メモリプロファイル

```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## CI/CD

### GitHub Actions (予定)

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.23'
      - run: make deps
      - run: make test
      - run: make build
```

## リリース

### バージョニング

Semantic Versioningに従う：`MAJOR.MINOR.PATCH`

### リリース手順

1. バージョンタグを作成
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. リリースノートを作成

3. バイナリをビルド
   ```bash
   make build
   ```

## 参考資料

- [Go Documentation](https://golang.org/doc/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Slack API Documentation](https://api.slack.com/)
