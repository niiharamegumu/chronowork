# ChronoWork

ChronoWork は、作業時間を追跡・管理するためのターミナルベースのTUIアプリケーションです。

## 特徴

- **TUIインターフェース**: ターミナル上で動作する直感的なユーザーインターフェース
- **時間追跡**: 作業の開始・停止を簡単に記録
- **プロジェクト管理**: プロジェクトとタグで作業を分類
- **データエクスポート**: CSVフォーマットでデータをエクスポート
- **クリーンアーキテクチャ**: テスト可能で保守性の高い設計

## アーキテクチャ

このプロジェクトはクリーンアーキテクチャに基づいて設計されています：

```
internal/
├── domain/       # ドメインエンティティ（純粋なGo構造体）
├── repository/   # リポジトリインターフェースとGORM実装
│   └── mock/     # テスト用モック実装
└── usecase/      # ビジネスロジック

container/        # 依存性注入コンテナ
widgets/          # TUIウィジェット
```

## 必要要件

- Go 1.21以上
- SQLite3

## インストール

```bash
git clone https://github.com/niiharamegumu/chronowork.git
cd chronowork
go build -o chronowork main.go
```

## 使い方

### 環境変数の設定

```bash
export CHRONOWORK_ROOT_PATH=/path/to/your/data
export DATABASE_NAME=sqlite.db  # オプション（デフォルト: sqlite.db）
```

### アプリケーションの起動

```bash
# 開発環境
make run

# 本番環境
make run-prod

# ビルド
make build
```

### キーバインディング

#### メインメニュー
- `w` - 作業一覧
- `p` - プロジェクト管理
- `t` - タグ管理
- `e` - データエクスポート
- `s` - 設定
- `q` - 終了
- `Esc` - メニューに戻る

#### 作業一覧
- `Enter` - 作業の追跡開始/停止
- `a` - 新規作業追加
- `u` - 作業編集
- `r` - 作業時間のリセット
- `d` - 作業削除
- `c` - 作業の確認状態切り替え
- `t` - タイトルをクリップボードにコピー
- `h` - 作業時間をクリップボードにコピー
- `s` - テーブルの先頭に移動
- `e` - テーブルの末尾に移動

#### プロジェクト/タグ管理
- `a` - 新規追加
- `u` - 編集
- `d` - 削除

## テスト

```bash
# 全テストの実行
go test -v ./...

# カバレッジ付きテスト
go test -cover ./...

# 特定パッケージのテスト
go test -v ./internal/usecase/...
```

## 開発

### ディレクトリ構造

```
.
├── app/                    # アプリケーション初期化
├── container/              # DIコンテナ
├── db/                     # データベース接続
├── internal/
│   ├── domain/            # ドメインエンティティ
│   ├── repository/        # リポジトリ層
│   │   └── mock/          # モック実装
│   └── usecase/           # ユースケース層
├── models/                # GORMモデル
├── service/               # TUIサービス
├── util/                  # ユーティリティ
│   ├── strutil/
│   └── timeutil/
├── widgets/               # TUIウィジェット
├── main.go
├── Makefile
└── go.mod
```

### 依存関係

- [tview](https://github.com/rivo/tview) - TUIフレームワーク
- [tcell](https://github.com/gdamore/tcell) - ターミナルセルベースのUI
- [GORM](https://gorm.io/) - ORMライブラリ
- [clipboard](https://golang.design/x/clipboard) - クリップボード操作

## ライセンス

MIT License

## 作者

Megumu Niihara
