package route

import (
	"himaplus-authn/view"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// エンドポイントのルーティング？

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
	// アクセストークンを返す
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error { // pbインスタンスのOnServe()フックを使って処理を鯖起動時にトリガーする // 処理はBindFunc()に渡す
		se.Router.GET("/google/access_token", func(re *core.RequestEvent) error { // core.RequestEventはreq, resを操作するためのメソッドを持つ
			// TODO: 色々な処理

			googleAccessToken := "token"

			// 値をJSON形式で返却
			return re.JSON(http.StatusOK, map[string]string{
				"googleAccessToken": googleAccessToken,
			})
		})
		// .Bind(apis.RequireAuth()) // HTTPメソッド関数にチェーンしてミドルウェアを追加できる

		return se.Next() // CONTEXT: OnServe()くんがエラーを
	})
}

// ルーティング
func Routing() *pocketbase.PocketBase {
	// pbインスタンス
	pb := pocketbase.New()

	// ルーティング
	fileRouting(pb)     // ファイル
	endpointRouting(pb) // カスタムエンドポイントを拡張

	return pb
}
