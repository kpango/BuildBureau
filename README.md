# BuildBureau

多層AIエージェント実装システム - 社長から平社員まで階層型マルチエージェント構成

## 概要

BuildBureauは、社長→部長→課長→平社員の階層型マルチエージェント構成を取るAIシステムです。各階層には秘書エージェントが存在し、上位エージェントからの指示を詳細化し記録・補佐・スケジューリングします。

### アーキテクチャ

```
クライアント
    ↓
社長エージェント + 社長秘書
    ↓
部長エージェント + 部長秘書  
    ↓
課長エージェント + 課長秘書
    ↓
平社員エージェント
```

### 主な特徴

- **階層型エージェント構造**: 社長、部長、課長、平社員の4層構造
- **秘書エージェント**: 各階層に秘書エージェントが存在し、タスク管理をサポート
- **gRPC通信**: エージェント間はgRPCで疎結合に通信
- **YAML設定**: すべての設定をYAMLファイルで管理
- **Slack通知**: 重要なイベントをSlackに自動通知
- **Terminal UI**: Bubble Teaによる対話型ターミナルUI
- **単一バイナリ**: Goで実装された単一バイナリで動作

## 技術スタック

- **言語**: Go 1.23+
- **AIエージェント**: Google ADK (Agent Development Kit) for Go
- **通信**: gRPC (Protocol Buffers)
- **UI**: Charmbracelet Bubble Tea
- **通知**: Slack API (slack-go)
- **設定**: YAML (gopkg.in/yaml.v3)

## インストール

### 前提条件

- Go 1.23以上
- protoc (Protocol Buffers compiler)

### ビルド

```bash
# 依存関係のインストール
make deps

# プロトコルバッファのコード生成（必要な場合）
make install-tools
make proto

# ビルド
make build
```

## 設定

`config.yaml`ファイルで全ての設定を管理します。

### 主要設定項目

#### エージェント設定

各エージェントタイプごとに以下を設定:

- `count`: エージェント数
- `model`: 使用するLLMモデル
- `instruction`: エージェントへのシステムプロンプト
- `allowTools`: ツール使用の許可
- `tools`: 使用可能なツールのリスト
- `timeout`: タイムアウト時間（秒）
- `retryCount`: リトライ回数

```yaml
agents:
  president:
    count: 1
    model: "gemini-2.0-flash-exp"
    instruction: |
      あなたは社長としてプロジェクト全体を俯瞰し方針を決定する立場です。
    allowTools: true
    tools:
      - web_search
      - knowledge_base
    timeout: 120
    retryCount: 3
```

#### Slack通知設定

```yaml
slack:
  enabled: true
  token: "${SLACK_BOT_TOKEN}"
  channelID: "${SLACK_CHANNEL_ID}"
  notifications:
    projectStart:
      enabled: true
      message: "🚀 プロジェクト「{{.ProjectName}}」が開始されました"
```

環境変数でトークンとチャンネルIDを設定:

```bash
export SLACK_BOT_TOKEN="xoxb-your-token"
export SLACK_CHANNEL_ID="C01234567"
```

#### UI設定

```yaml
ui:
  enableTUI: true
  refreshRate: 100  # ミリ秒
  theme: "default"
  showProgress: true
  logLevel: "info"
```

## 使い方

### 基本的な実行

```bash
# デフォルト設定で実行
./bin/buildbureau

# カスタム設定ファイルを指定
CONFIG_PATH=/path/to/config.yaml ./bin/buildbureau
```

### Terminal UI

TUIが有効な場合、対話型のターミナルインターフェースが起動します:

- プロジェクト要件を入力
- `Alt+Enter`: 要件を送信してプロジェクト開始
- `Esc`: 終了

### エージェントの動作フロー

1. **社長エージェント**: クライアントからの要件を受け取り、全体計画を立案
2. **社長秘書**: 要件を記録し、詳細化して部長秘書へ
3. **部長エージェント**: タスクを課長単位に分割
4. **部長秘書**: タスクを詳細化し、課長秘書へ
5. **課長エージェント**: 実装計画と仕様書を策定
6. **課長秘書**: 実装手順のドラフトを作成
7. **平社員エージェント**: 具体的な実装を実行

## 開発

### ディレクトリ構造

```
BuildBureau/
├── cmd/
│   └── buildbureau/      # メインアプリケーション
│       └── main.go
├── internal/
│   ├── agent/            # エージェント実装
│   ├── config/           # 設定管理
│   ├── grpc/             # gRPCサービス実装
│   ├── slack/            # Slack通知
│   └── ui/               # Terminal UI
├── proto/
│   └── buildbureau/v1/   # Protocol Buffers定義
├── pkg/                  # 公開パッケージ
├── config.yaml           # デフォルト設定
├── Makefile             # ビルドスクリプト
└── go.mod               # Go依存関係
```

### テスト

```bash
make test
```

### フォーマットとLint

```bash
make lint
```

## gRPCサービス

各階層でgRPCサービスを定義:

- **PresidentService**: プロジェクト計画立案
- **DepartmentManagerService**: タスク分割
- **SectionManagerService**: 実装計画策定
- **EmployeeService**: タスク実行

詳細は`proto/buildbureau/v1/service.proto`を参照。

## Slack通知

以下のイベントでSlack通知が送信されます:

- プロジェクト開始
- タスク完了
- エラー発生
- プロジェクト完了

通知の有効/無効や内容は`config.yaml`で設定可能。

## ライセンス

このプロジェクトのライセンスについては[LICENSE](LICENSE)ファイルを参照してください。

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。

## TODO

- [ ] Google ADK統合の実装
- [ ] gRPCサービスの完全実装
- [ ] エージェント間通信の実装
- [ ] ナレッジベースの実装
- [ ] ツールシステムの実装
- [ ] ストリーミング対応
- [ ] エラーハンドリングの強化
- [ ] テストカバレッジの向上
- [ ] ドキュメントの充実