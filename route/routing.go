package route

import (
	"fmt"
	"himaplus-authn/view"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// グローバルミドルウェア
func globalMiddleware(pb *pocketbase.PocketBase) {
	// 現状なし
}

// ファイルのルーティング
func fileRouting(pb *pocketbase.PocketBase) {
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error { // CONTEXT: OnServe.BindFuncで鯖起動時のイベントの関数を設定
		// テスト用ページを表示

		// 静的ファイルのバイナリへの埋め込み
		indexStaticFS, err := view.EmbedIndexFile()
		if err != nil {
			return err
		}
		// CONTEXT: se.Routerはpbが内部的に利用しているHTTPルーターで、ここにエンドポイントを追加できる
		// ルーティング
		se.Router.GET("/{path...}", apis.Static(indexStaticFS, false)) // 第三引数はキャッシュの有無 // "/index/{path...}"

		// 任意のバスに対して適切な静的コンテンツを提供

		// 静的ファイルの埋め込み
		publicStaticFS, err := view.EmbedStaticFile()
		if err != nil {
			return err
		}
		// ルーティング
		se.Router.GET("/public/{path...}", apis.Static(publicStaticFS, false))

		return se.Next()
	})
}

// カスタムエンドポイント
func endpointRouting(pb *pocketbase.PocketBase) {
	// TODO: 認証状態を確認

	// アクセストークンを返す
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error { // pbインスタンスのOnServe()フックを使って処理を鯖起動時にトリガーする // 処理はBindFunc()に渡す
		se.Router.GET("/google/access_token", func(re *core.RequestEvent) error { // core.RequestEventはreq, resを操作するためのメソッドを持つ
			// TODO: 色々な処理
			// 認証済みユーザー（:アクセスを許されたユーザー）のレコードを取得authUserRecord
			userRecord := re.Auth
			fmt.Printf("authUserRecord: %v\n", userRecord)

			// アクセストークンを取得

			// OnRecordAuthRefreshRequestと併用する？

			googleAccessToken := "token"

			// 値をJSON形式で返却
			return re.JSON(http.StatusOK, map[string]string{
				"googleAccessToken": googleAccessToken,
			})
		}).Bind(apis.RequireAuth("_superusers", "users")) // HTTPメソッド関数にチェーンしてミドルウェアを追加できる

		return se.Next()
	})

	// リフレッシュの検証
	// fires only for "users" and "managers" auth collections
	pb.OnRecordAuthRefreshRequest("users").BindFunc(func(re *core.RecordAuthRefreshRequestEvent) error {
		// // e.App
		// // e.Collection
		// // e.Record
		// // and all RequestEvent fields...
		// fmt.Printf("re.App: %v\n", re.App)
		// fmt.Printf("re.Collection: %v\n", re.Collection)
		// fmt.Printf("re.Record: %v\n", re.Record)
		// fmt.Printf("re.Collection.OAuth2: %v\n", re.Collection.OAuth2)
		// fmt.Printf("re.Auth.Collection().OAuth2: %v\n", re.Auth.Collection().OAuth2)
		// fmt.Printf("re.Record.Collection().OAuth2: %v\n", re.Record.Collection().OAuth2)

		// fmt.Println()
		// fmt.Println()

		// re.App.OnRecordAuthWithOAuth2Request().BindFunc(func(re *core.RecordAuthWithOAuth2RequestEvent) error {
		// 	fmt.Printf("re.OAuth2User.AccessToken: %v\n", re.OAuth2User.AccessToken)
		// 	return re.Next()
		// })

		return re.Next()
	})

	// // アクセストークン取得検証
	// pb.OnRecordAuthWithOAuth2Request("users").BindFunc(func(re *core.RecordAuthWithOAuth2RequestEvent) error {
	// 	fmt.Println("認証後")

	// 	fmt.Printf("re.OAuth2User.AccessToken: %v\n", re.OAuth2User.AccessToken)
	// 	fmt.Printf("re.OAuth2User.RefreshToken: %v\n", re.OAuth2User.RefreshToken) // access_type=offlineを指定してないとonlineとみなされてもらえない

	// 	return re.Next()
	// })

	// 認証リクエスト前の調整検証
	// // 認証リクエスト前にURLに、refresh token取得に必要な情報を追加
	// pb.OnRecordAuthWithOAuth2Request("users").BindFunc(func(re *core.RecordAuthWithOAuth2RequestEvent) error {
	// 	// before処理

	// 	// 認証URLをurlとしてパースして取得
	// 	authUrl, err := url.Parse(re.ProviderClient.AuthURL())
	// 	fmt.Printf("authUrl: %v\n", authUrl)
	// 	fmt.Printf("authUrl.RawQuery: %v\n", authUrl.RawQuery)
	// 	fmt.Printf("re.ProviderClient.Scopes(): %v\n", re.ProviderClient.Scopes())
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// URLに必要なクエパラを追加

	// 	queryPara := authUrl.Query() // クエリを取得
	// 	fmt.Printf("queryPara: %v\n", queryPara)

	// 	// 値をセット
	// 	queryPara.Set("access_type", "offline")
	// 	queryPara.Set("prompt", "consent")
	// 	queryPara.Set("redirect_uri", re.ProviderClient.RedirectURL())
	// 	queryPara.Set("scope", strings.Join([]string{
	// 		"https://www.googleapis.com/auth/calendar",
	// 		"https://www.googleapis.com/auth/calendar.readonly",
	// 		"https://www.googleapis.com/auth/userinfo.email",
	// 		"https://www.googleapis.com/auth/userinfo.profile",
	// 	}, " "))

	// 	authUrl.RawQuery = queryPara.Encode() // 元のクエパラに上書き
	// 	fmt.Printf("authUrl: %v\n", authUrl)

	// 	// 更新したURLを適用
	// 	re.ProviderClient.SetAuthURL(authUrl.String())
	// 	fmt.Printf("re.ProviderClient.AuthURL(): %v\n", re.ProviderClient.AuthURL())

	// 	return re.Next()
	// })

	// TODO: 認証リクエスト後にrefresh tokenを保存
	pb.OnRecordAuthWithOAuth2Request("users").BindFunc(func(re *core.RecordAuthWithOAuth2RequestEvent) error {
		if err := re.Next(); err != nil {
			return err
		}

		// after処理

		fmt.Println()

		fmt.Printf("re.ProviderClient.AuthURL(): %v\n", re.ProviderClient.AuthURL())
		fmt.Printf("re.Auth: %v\n", re.Auth)
		fmt.Printf("re: %v\n", re)
		fmt.Printf("re.OAuth2User.AccessToken: %v\n", re.OAuth2User.AccessToken)
		fmt.Printf("re.OAuth2User.RefreshToken: %v\n", re.OAuth2User.RefreshToken) // access_type=offlineを指定してないとonlineとみなされてもらえない

		return nil
	})

	// HACK: 前後処理をまとめれる？
}

// ルーティング
func Routing() *pocketbase.PocketBase {
	// pbインスタンス
	pb := pocketbase.New()

	// 拡張

	globalMiddleware(pb) // ミドルウェア
	fileRouting(pb)      // ファイル公開
	endpointRouting(pb)  // カスタムエンドポイントのルーティング

	return pb
}
