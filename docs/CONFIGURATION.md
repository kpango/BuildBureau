# Configuration Guide

BuildBureauの詳細な設定ガイド

## 設定ファイルの構造

BuildBureauは`config.yaml`で全ての設定を管理します。

## エージェント設定 (agents)

各エージェント種別ごとに以下のパラメータを設定できます：

### 共通パラメータ

```yaml
agents:
  <agent_type>:
    count: <number>           # エージェント数
    model: <string>           # 使用するLLMモデル名
    instruction: <string>     # システムプロンプト
    allowTools: <boolean>     # ツール使用許可
    tools: [<strings>]        # 使用可能なツールリスト
    timeout: <seconds>        # タイムアウト時間
    retryCount: <number>      # リトライ回数
```

### エージェント種別

#### 1. president (社長)
```yaml
president:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    あなたは社長としてプロジェクト全体を俯瞰し方針を決定する立場です。
    クライアントからの要求を理解し、プロジェクト全体の計画を立案してください。
  allowTools: true
  tools:
    - web_search      # Web検索
    - knowledge_base  # ナレッジベースアクセス
  timeout: 120
  retryCount: 3
```

#### 2. president_secretary (社長秘書)
```yaml
president_secretary:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    あなたは社長の秘書です。社長の指示を受けて要件を記録し、
    社内ナレッジベースを更新してください。
  allowTools: true
  tools:
    - knowledge_base
    - document_manager
  timeout: 60
  retryCount: 3
```

#### 3. department_manager (部長)
```yaml
department_manager:
  count: 1
  model: "gemini-2.0-flash-exp"
  instruction: |
    あなたは部長としてプロジェクト全体を課長単位に分割する責任者です。
  allowTools: true
  tools:
    - web_search
    - knowledge_base
  timeout: 120
  retryCount: 3
```

#### 4. section_manager (課長)
```yaml
section_manager:
  count: 3  # 複数人配置可能
  model: "gemini-2.0-flash-exp"
  instruction: |
    あなたは課長として詳細な実装計画と最終仕様書を策定する責任者です。
  allowTools: true
  tools:
    - code_analyzer
    - knowledge_base
  timeout: 90
  retryCount: 3
```

#### 5. employee (平社員)
```yaml
employee:
  count: 6  # 複数人配置可能
  model: "gemini-2.0-flash-exp"
  instruction: |
    あなたは与えられた仕様に基づき実装を行うエンジニアです。
  allowTools: true
  tools:
    - code_execution
    - file_operations
    - knowledge_base
  timeout: 180
  retryCount: 3
```

## LLM設定 (llm)

```yaml
llm:
  provider: "google"                                    # プロバイダー名
  apiEndpoint: "https://generativelanguage.googleapis.com"  # APIエンドポイント
  defaultModel: "gemini-2.0-flash-exp"                 # デフォルトモデル
  maxTokens: 8192                                       # 最大トークン数
  temperature: 0.7                                      # 温度パラメータ (0.0-1.0)
  topP: 0.95                                           # Top-Pサンプリング
```

### プロバイダー設定

現在サポート予定のプロバイダー：
- `google`: Google AI (Gemini)
- `openai`: OpenAI (GPT-4など)
- `anthropic`: Anthropic (Claude)

### モデル選択

推奨モデル：
- 高速処理: `gemini-2.0-flash-exp`
- 高品質: `gemini-2.5-pro`
- バランス: `gemini-2.0-flash-exp`

## gRPC設定 (grpc)

```yaml
grpc:
  port: 50051                  # gRPCサーバーポート
  maxMessageSize: 10485760     # 最大メッセージサイズ (バイト)
  timeout: 300                 # タイムアウト (秒)
  enableReflection: true       # リフレクション有効化
```

### ポート設定

- デフォルト: `50051`
- ファイアウォールでこのポートを開放する必要がある場合があります

### メッセージサイズ

- デフォルト: 10MB
- 大きなファイルを扱う場合は増やす

## Slack通知設定 (slack)

```yaml
slack:
  enabled: true                      # Slack通知の有効化
  token: "${SLACK_BOT_TOKEN}"        # Botトークン (環境変数)
  channelID: "${SLACK_CHANNEL_ID}"   # チャンネルID (環境変数)
  retryCount: 3                      # リトライ回数
  timeout: 10                        # タイムアウト (秒)
  
  notifications:
    projectStart:
      enabled: true
      message: "🚀 プロジェクト「{{.ProjectName}}」が開始されました"
    
    taskComplete:
      enabled: true
      message: "✅ タスク「{{.TaskName}}」が完了しました ({{.Agent}})"
    
    error:
      enabled: true
      message: "❌ エラーが発生しました: {{.ErrorMessage}} ({{.Agent}})"
    
    projectComplete:
      enabled: true
      message: "🎉 プロジェクト「{{.ProjectName}}」が完了しました！"
```

### Slack Bot設定手順

1. [Slack API](https://api.slack.com/apps)でアプリを作成
2. Bot Token Scopesに以下を追加：
   - `chat:write`
   - `chat:write.public`
3. ワークスペースにインストール
4. Bot User OAuth Tokenを取得
5. 環境変数に設定：
   ```bash
   export SLACK_BOT_TOKEN="xoxb-your-token"
   export SLACK_CHANNEL_ID="C01234567"
   ```

### メッセージテンプレート

利用可能な変数：
- `{{.ProjectName}}`: プロジェクト名
- `{{.TaskName}}`: タスク名
- `{{.Agent}}`: エージェントID
- `{{.ErrorMessage}}`: エラーメッセージ
- `{{.Timestamp}}`: タイムスタンプ

## UI設定 (ui)

```yaml
ui:
  enableTUI: true        # Terminal UIの有効化
  refreshRate: 100       # 更新間隔 (ミリ秒)
  theme: "default"       # テーマ
  showProgress: true     # プログレス表示
  logLevel: "info"       # ログレベル
```

### ログレベル

- `debug`: デバッグ情報を含む全てのログ
- `info`: 通常の情報ログ
- `warn`: 警告のみ
- `error`: エラーのみ

### テーマ

現在利用可能なテーマ：
- `default`: デフォルトテーマ

## システム設定 (system)

```yaml
system:
  workDir: "./work"              # 作業ディレクトリ
  logDir: "./logs"               # ログディレクトリ
  cacheDir: "./cache"            # キャッシュディレクトリ
  maxConcurrentTasks: 10         # 同時実行タスク数
  globalTimeout: 3600            # グローバルタイムアウト (秒)
```

### ディレクトリ構造

```
BuildBureau/
├── work/      # 作業用一時ファイル
├── logs/      # ログファイル
└── cache/     # キャッシュファイル
```

## 環境変数

### 必須環境変数

Slack通知を使用する場合：
```bash
export SLACK_BOT_TOKEN="xoxb-..."
export SLACK_CHANNEL_ID="C..."
```

Google AI APIを使用する場合：
```bash
export GOOGLE_AI_API_KEY="..."
```

### オプション環境変数

```bash
# カスタム設定ファイルパス
export CONFIG_PATH="/path/to/custom/config.yaml"

# ログレベルの上書き
export LOG_LEVEL="debug"
```

## 設定例

### 開発環境用設定

```yaml
agents:
  president:
    count: 1
    timeout: 60
  # ... 他のエージェント（タイムアウトを短く）

slack:
  enabled: false  # 開発時は通知無効化

ui:
  enableTUI: true
  logLevel: "debug"  # デバッグログ有効
```

### 本番環境用設定

```yaml
agents:
  president:
    count: 1
    timeout: 180
  section_manager:
    count: 5  # スケールアップ
  employee:
    count: 20  # スケールアップ

slack:
  enabled: true  # 通知有効化

system:
  maxConcurrentTasks: 20  # 並列度向上

ui:
  logLevel: "info"  # 情報ログのみ
```

### 高負荷環境用設定

```yaml
grpc:
  maxMessageSize: 52428800  # 50MB

system:
  maxConcurrentTasks: 50
  globalTimeout: 7200  # 2時間

agents:
  employee:
    count: 50
    timeout: 300
```

## トラブルシューティング

### 設定ファイルのバリデーション

設定ファイルの構文チェック：
```bash
# YAML構文チェック
yamllint config.yaml

# BuildBureauで検証
./bin/buildbureau --validate-config  # (未実装)
```

### よくあるエラー

1. **"failed to load config"**
   - YAMLの構文エラーをチェック
   - インデントが正しいか確認

2. **"Slack token is required"**
   - 環境変数 `SLACK_BOT_TOKEN` を設定
   - `slack.enabled: false` にする

3. **"president agent count must be at least 1"**
   - 必須エージェントのカウントを確認

## ベストプラクティス

1. **機密情報の管理**
   - トークンは環境変数で管理
   - `.env`ファイルを`.gitignore`に追加

2. **タイムアウト設定**
   - エージェントの役割に応じて適切に設定
   - 実装を伴う作業は長めに設定

3. **リトライ回数**
   - ネットワーク不安定な環境では多めに設定
   - 無限ループを避けるため上限を設ける

4. **ログレベル**
   - 開発時は`debug`
   - 本番時は`info`または`warn`
