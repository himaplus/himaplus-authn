package main

import (
	"himaplus-authn/view"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	// pbインスタンス
	app := pocketbase.New()

	// イベント追加

	// テスト用ページを表示
	app.OnServe().BindFunc(func(se *core.ServeEvent) error { // CONTEXT: OnServe.BindFuncで鯖起動時のイベントの関数を設定
		staticFS, err := view.EmbedIndexFile()
		if err != nil {
			return err
		}
		// CONTEXT: se.Routerはpbが内部的に利用しているHTTPルーターで、ここにエンドポイントを追加できる
		se.Router.GET("/{path...}", apis.Static(staticFS, false)) // 第三引数はキャッシュの有無 // "/index/{path...}"
		return se.Next()
	})

	// 任意のバスに対して適切な静的コンテンツを提供
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		staticFS, err := view.EmbedStaticFile()
		if err != nil {
			return err
		}
		se.Router.GET("/public/{path...}", apis.Static(staticFS, false))
		return se.Next()
	})

	// 鯖起動
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
