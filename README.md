# himaplus-authn

## 概要

ひまぷらの認証サーバ  
go-echo-pocketbaseの認証サーバ  

### 環境

Visual Studio Code: 1.88.1  
golang.Go: v0.41.4  
image Golang: go version go1.22.2 linux/amd64
TODO: version
echo: v
pocketbase: v

## 環境構築

[docker-himaplus](https://github.com/unSerori/docker-himaplus)を使ってDokcerコンテナーで開発・デプロイする形を想定している  
インストール手順は[docker-himaplusのインストール項目](https://github.com/unSerori/docker-himaplus/blob/main/README.md#インストール)に記載  
cloneしてスクリプト実行で、自動的にコンテナー作成と開発環境（: またはデプロイ）を行う  

## API仕様書

エンドポイント、リクエストレスポンスの形式、その他情報のAPIの仕様書。

### エンドポインツ

TODO: ここにエンドポイント仕様書

### API仕様書てんぷれ

<details>
  <summary>＊○○＊するエンドポイント</summary>

- **URL:** `/＊エンドポイントパス＊`
- **メソッド:** ＊HTTPメソッド名＊
- **説明:** ＊○○＊
- **リクエスト:**
  - ヘッダー:
    - `＊HTTPヘッダー名＊`: ＊HTTPヘッダー値＊
  - ボディ:
    ＊さまざまな形式のボディ値＊

- **レスポンス:**
  - ステータスコード: ＊ステータスコード ステータスメッセージ＊
    - ボディ:
      ＊さまざまな形式のレスポンスデータ（基本はJSON）＊

      ```json
      {
        "srvResMsg":  "レスポンスステータスメッセージ",
        "srvResData": {
        
        },
      }
      ```

</details>

## .ENV

.evnファイルの各項目と説明

```env:.env
```
