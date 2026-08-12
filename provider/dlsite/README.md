# DLSite Provider

DLSiteプロバイダーは、[DLSite](https://www.dlsite.com/)から動画・ボイス作品のメタデータを取得するための[metatube-sdk-go](https://github.com/metatube-community/metatube-sdk-go)プロバイダー実装です。

## 概要

このプロバイダーは、DLSiteで販売されている**動画・ボイス作品（RJ番号）のみ**を対象としています。作品ID検索、キーワード検索、および詳細メタデータ取得をサポートします。

### 対象コンテンツ

- ✅ **動画作品** (RJ番号)
- ✅ **ボイス作品** (RJ番号、主に同人音声作品)

### 対象外コンテンツ

- ❌ **ゲーム作品** (VJ番号)
- ❌ **コミック・マンガ作品** (BJ番号)
- ❌ その他のコンテンツタイプ

## 機能

- **作品ID正規化**: RJ番号の抽出と正規化（大文字変換、パターンマッチング）
- **URL解析**: DLSite URLから作品IDを抽出
- **キーワード検索**: 作品タイトルやキーワードによる検索（RJ作品のみ）
- **詳細メタデータ取得**: タイトル、サークル名、販売日、あらすじ、画像、ジャンル等
- **年齢確認の自動処理**: 成人向けコンテンツへの自動アクセス
- **コンテンツタイプフィルタリング**: 動画・ボイス作品のみを抽出

## インストール

```bash
go get github.com/metatube-community/metatube-sdk-go
```

## 基本的な使用例

### プロバイダーの初期化

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/metatube-community/metatube-sdk-go/provider/dlsite"
)

func main() {
    // プロバイダーのインスタンスを作成
    provider := dlsite.New()
    
    // プロバイダー情報の取得
    fmt.Printf("Provider: %s\n", provider.Name())        // "dlsite"
    fmt.Printf("URL: %s\n", provider.URL().String())     // "https://www.dlsite.com"
    fmt.Printf("Language: %s\n", provider.Language())    // "ja"
}
```

### 作品IDによるメタデータ取得

```go
// 作品IDで詳細情報を取得
info, err := provider.GetMovieInfoByID("RJ123456")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("ID: %s\n", info.ID)
fmt.Printf("Title: %s\n", info.Title)
fmt.Printf("Maker: %s\n", info.Maker)
fmt.Printf("Release Date: %s\n", info.ReleaseDate)
fmt.Printf("Summary: %s\n", info.Summary)
fmt.Printf("Cover URL: %s\n", info.CoverURL)
fmt.Printf("Genres: %v\n", info.Genres)
```

### URLによるメタデータ取得

```go
// DLSite URLから直接メタデータを取得
url := "https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html"
info, err := provider.GetMovieInfoByURL(url)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Title: %s\n", info.Title)
```

### キーワード検索

```go
// キーワードで作品を検索
results, err := provider.SearchMovie("ボイスドラマ")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d results\n", len(results))
for _, result := range results {
    fmt.Printf("- %s: %s\n", result.ID, result.Title)
    fmt.Printf("  URL: %s\n", result.Homepage)
    fmt.Printf("  Thumbnail: %s\n", result.ThumbURL)
}
```

### 作品IDによる直接検索

```go
// 作品IDを含むキーワードで検索すると、その作品の詳細情報を返す
results, err := provider.SearchMovie("RJ123456")
if err != nil {
    log.Fatal(err)
}

if len(results) > 0 {
    fmt.Printf("Found: %s\n", results[0].Title)
}
```

## 作品ID形式

DLSiteの作品IDは以下の形式です：

```
RJ + 6〜8桁の数字
```

### 有効な作品ID例

- `RJ123456` (6桁)
- `RJ01234567` (7桁)
- `RJ12345678` (8桁)

### 無効な作品ID例（対象外）

- `VJ012345` (ゲーム作品 - 非対応)
- `BJ234567` (コミック作品 - 非対応)
- `RE123456` (その他のコンテンツ - 非対応)

## 取得可能なメタデータ

`GetMovieInfoByID()` および `GetMovieInfoByURL()` で取得できる情報：

| フィールド | 説明 | 例 |
|----------|------|-----|
| `ID` | 作品ID | `RJ123456` |
| `Number` | 作品番号（IDと同じ） | `RJ123456` |
| `Title` | 作品タイトル | `サンプルボイス作品` |
| `Maker` | サークル名 | `サンプルサークル` |
| `ReleaseDate` | 販売日 | `2023-12-25` |
| `Summary` | あらすじ・説明 | `作品の詳細説明...` |
| `CoverURL` | カバー画像URL | `https://...` |
| `ThumbURL` | サムネイル画像URL | `https://...` |
| `Genres` | ジャンル・タグ配列 | `["ボイスドラマ", "癒し"]` |
| `Provider` | プロバイダー名 | `dlsite` |
| `Homepage` | 作品ページURL | `https://www.dlsite.com/...` |

## エラーハンドリング

プロバイダーは以下のエラータイプを返します：

```go
var (
    ErrInvalidWorkID          // 無効な作品ID形式
    ErrUnsupportedContentType // 非対応のコンテンツタイプ（VJ、BJ等）
    ErrWorkNotFound           // 作品が見つからない（HTTP 404）
    ErrParseError             // HTML解析エラー
    ErrNetworkError           // ネットワークエラー
    ErrValidationError        // メタデータ検証エラー
)
```

### エラー処理の例

```go
info, err := provider.GetMovieInfoByID("VJ012345")
if err != nil {
    if errors.Is(err, dlsite.ErrUnsupportedContentType) {
        fmt.Println("このコンテンツタイプは対応していません")
    } else if errors.Is(err, dlsite.ErrWorkNotFound) {
        fmt.Println("作品が見つかりません")
    } else {
        log.Fatal(err)
    }
}
```

## 制限事項

1. **コンテンツタイプ**: RJ番号の動画・ボイス作品のみをサポート
2. **言語**: 日本語のみ（DLSiteの日本語ページから情報を取得）
3. **年齢制限**: 成人向けコンテンツへのアクセスには年齢確認Cookieを自動設定
4. **レート制限**: DLSiteのレート制限に従う必要があります
5. **HTML構造依存**: DLSiteのHTML構造変更により動作しなくなる可能性があります

## 高度な使用例

### タイムアウトの設定

```go
provider := dlsite.New()

// リクエストタイムアウトを60秒に設定
provider.SetRequestTimeout(60 * time.Second)
```

### 作品IDの正規化

```go
provider := dlsite.New()

// 様々な形式の入力から作品IDを抽出
id1 := provider.NormalizeMovieID("rj123456")                    // "RJ123456"
id2 := provider.NormalizeMovieID("RJ123456")                    // "RJ123456"
id3 := provider.NormalizeMovieID("https://dlsite.com/RJ123456") // "RJ123456"
id4 := provider.NormalizeMovieID("VJ012345")                    // "" (非対応)
id5 := provider.NormalizeMovieID("invalid")                     // "" (無効)

if id1 == "" {
    fmt.Println("無効な作品IDまたは非対応のコンテンツタイプ")
}
```

### URLからの作品ID抽出

```go
provider := dlsite.New()

url := "https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html"
id, err := provider.ParseMovieIDFromURL(url)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Extracted ID: %s\n", id) // "RJ123456"
```

## テスト

```bash
# 全テストを実行
go test ./provider/dlsite/...

# 統合テストを除外（ネットワークアクセスなし）
go test -short ./provider/dlsite/...

# Property-Based Testsを含む全テストを実行
go test -v ./provider/dlsite/...
```

## ライセンス

このプロバイダーは metatube-sdk-go プロジェクトの一部であり、同じライセンスの下で提供されます。

## 貢献

バグ報告や機能リクエストは、[metatube-sdk-go リポジトリ](https://github.com/metatube-community/metatube-sdk-go)のIssuesセクションにお願いします。

## 関連リンク

- [DLSite 公式サイト](https://www.dlsite.com/)
- [metatube-sdk-go](https://github.com/metatube-community/metatube-sdk-go)
- [GoDoc](https://pkg.go.dev/github.com/metatube-community/metatube-sdk-go/provider/dlsite)
