# Quick Start Guide

BuildBureauを5分で始める

## 1. インストール

### 前提条件

- Go 1.23以上

### ビルド

```bash
git clone https://github.com/kpango/BuildBureau.git
cd BuildBureau
make deps
make build
```

## 2. 設定

### 最小設定で開始

`config.yaml`はそのまま使用可能です。

### Slack通知を有効にする（オプション）

```bash
# .envファイルを作成
cp .env.example .env

# .envを編集してトークンを設定
export SLACK_BOT_TOKEN="xoxb-your-token"
export SLACK_CHANNEL_ID="C01234567"
```

Slack通知が不要な場合：

```yaml
# config.yamlで無効化
slack:
  enabled: false
```

## 3. 実行

### デフォルト設定で実行

```bash
./bin/buildbureau
```

### Terminal UIの操作

起動すると対話型のターミナルUIが表示されます：

```
🏢 BuildBureau - マルチレイヤー AI エージェントシステム

要件入力:
┌──────────────────────────────────────┐
│ プロジェクトの要件を入力してください...│
│                                      │
│                                      │
└──────────────────────────────────────┘

Alt+Enter: 送信 | Esc: 終了
```

### プロジェクト要件の入力

1. テキストエリアにプロジェクトの要件を入力
2. `Alt+Enter`を押して送信
3. エージェントが処理を開始

## 4. 動作確認

### 例：シンプルなプロジェクト

```
Webサイトの問い合わせフォームを作成してください。
以下の機能が必要です：
- 名前、メールアドレス、メッセージの入力欄
- バリデーション
- 送信確認
```

### エージェントの動作

1. 社長エージェントが要件を分析
2. 部長エージェントがタスクを分割
3. 課長エージェントが実装計画を作成
4. 平社員エージェントが実装を実行

## 5. 設定のカスタマイズ

### エージェント数の変更

```yaml
# config.yaml
agents:
  section_manager:
    count: 5  # 課長を5人に増やす
  employee:
    count: 20  # 平社員を20人に増やす
```

### タイムアウトの調整

```yaml
agents:
  president:
    timeout: 180  # 180秒に延長
```

### LLMモデルの変更

```yaml
agents:
  president:
    model: "gemini-2.5-pro"  # より高性能なモデルに
```

## トラブルシューティング

### ビルドエラー

```bash
make clean
make deps
make build
```

### 設定エラー

```bash
# YAMLの構文チェック
cat config.yaml | grep -E "^\s*-"
```

### Slack通知が届かない

1. トークンが正しいか確認
2. チャンネルIDが正しいか確認
3. Botがチャンネルに追加されているか確認

## 次のステップ

- [設定ガイド](docs/CONFIGURATION.md)で詳細な設定方法を確認
- [アーキテクチャドキュメント](docs/ARCHITECTURE.md)でシステムの仕組みを理解
- [開発ガイド](docs/DEVELOPMENT.md)でカスタマイズ方法を学習

## よくある質問

### Q: LLMの実装はどうなっていますか？

A: 現在のバージョンは基盤実装です。Google ADKとの統合は今後実装予定です。

### Q: エージェントはどこで動作しますか？

A: 現在は単一プロセス内で動作します。gRPCインターフェースにより、将来的に分散実行も可能です。

### Q: カスタムエージェントを追加できますか？

A: はい。[開発ガイド](docs/DEVELOPMENT.md)を参照してください。

### Q: 商用利用できますか？

A: ライセンスを確認してください。

## サポート

- バグ報告: [GitHub Issues](https://github.com/kpango/BuildBureau/issues)
- 質問: [GitHub Discussions](https://github.com/kpango/BuildBureau/discussions)
- 貢献: [Contributing Guide](CONTRIBUTING.md)
