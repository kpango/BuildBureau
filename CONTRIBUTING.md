# Contributing to BuildBureau

BuildBureauへの貢献ありがとうございます！

## 行動規範

プロジェクトに参加するすべての人は、尊重と礼儀を持って行動することが期待されます。

## 貢献の方法

### バグ報告

バグを見つけた場合：

1. 既存のissueを確認
2. 新しいissueを作成し、以下を含める：
   - バグの詳細な説明
   - 再現手順
   - 期待される動作
   - 実際の動作
   - 環境情報（OS、Goバージョン等）

### 機能提案

新機能を提案する場合：

1. issueで提案を議論
2. 実装前にメンテナーの承認を得る
3. 実装の詳細を含める

### プルリクエスト

1. **フォークとクローン**
   ```bash
   git clone https://github.com/YOUR_USERNAME/BuildBureau.git
   cd BuildBureau
   ```

2. **ブランチ作成**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **変更を実装**
   - コードを書く
   - テストを追加
   - ドキュメントを更新

4. **テスト**
   ```bash
   make test
   make lint
   ```

5. **コミット**
   ```bash
   git add .
   git commit -m "feat: Add your feature"
   ```

6. **プッシュ**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **プルリクエスト作成**
   - GitHubでPRを作成
   - 変更内容を説明
   - 関連issueをリンク

## コーディング規約

### Go スタイルガイド

- [Effective Go](https://golang.org/doc/effective_go.html)に従う
- `gofmt`でフォーマット
- `go vet`でチェック

### 命名規則

```go
// パッケージ: 小文字、単語
package agent

// エクスポート型: PascalCase
type AgentPool struct {}

// 非エクスポート型: camelCase
type agentInternal struct {}

// 関数: PascalCase (エクスポート), camelCase (非エクスポート)
func NewAgent() {}
func processTask() {}

// 定数: PascalCase (エクスポート), camelCase (非エクスポート)
const MaxRetryCount = 3
const defaultTimeout = 60
```

### コメント

```go
// Package agent provides AI agent implementations.
package agent

// Agent represents an AI agent interface.
// All agent types must implement this interface.
type Agent interface {
    // Process handles the given input and returns output.
    Process(ctx context.Context, input interface{}) (interface{}, error)
}
```

### エラーハンドリング

```go
// 良い例
if err != nil {
    return fmt.Errorf("failed to process task: %w", err)
}

// エラーは適切にラップする
if err := doSomething(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## テスト

### テストの必須事項

- すべての公開関数にテストを書く
- テストカバレッジ80%以上を目指す
- テーブル駆動テストを使用

### テストの例

```go
func TestNewAgent(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        want    string
        wantErr bool
    }{
        {"valid", "agent-1", "agent-1", false},
        {"empty id", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            agent, err := NewAgent(tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewAgent() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if agent != nil && agent.ID() != tt.want {
                t.Errorf("NewAgent() ID = %v, want %v", agent.ID(), tt.want)
            }
        })
    }
}
```

## ドキュメント

### コードドキュメント

- すべての公開APIにGoDocコメントを追加
- 複雑なロジックには説明コメントを追加

### README更新

機能追加時は以下を更新：

- README.md
- docs/ARCHITECTURE.md（必要に応じて）
- docs/CONFIGURATION.md（設定追加時）

## コミットメッセージ

### フォーマット

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（空白、フォーマット等）
- `refactor`: バグ修正も機能追加もしないコード変更
- `perf`: パフォーマンス改善
- `test`: テストの追加・修正
- `chore`: ビルドプロセスやツールの変更

### 例

```
feat(agent): Add support for streaming responses

Implement streaming support for long-running agent tasks.
This allows clients to receive progress updates in real-time.

Closes #123
```

## プルリクエストのレビュー

### レビュアー向け

- コードの品質を確認
- テストの十分性を確認
- ドキュメントの更新を確認
- 建設的なフィードバックを提供

### 著者向け

- フィードバックに対応
- 議論が必要な場合はコメントで説明
- 変更を素早く反映

## リリースプロセス

### バージョニング

Semantic Versioning (SemVer) を使用：

- `MAJOR`: 互換性のない変更
- `MINOR`: 後方互換性のある機能追加
- `PATCH`: 後方互換性のあるバグ修正

### リリース手順

1. CHANGELOGを更新
2. バージョンタグを作成
3. GitHubでリリースを作成
4. リリースノートを記述

## コミュニティ

### 質問

- GitHubのDiscussionsを使用
- issueでバグ報告

### コミュニケーション

- 日本語・英語どちらでもOK
- 丁寧で建設的なコミュニケーションを心がける

## ライセンス

貢献されたコードは、プロジェクトと同じライセンス（LICENSE参照）の下で提供されます。

## 謝辞

貢献者の皆様に感謝します！
