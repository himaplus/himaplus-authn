package view

import (
	"embed"
	"io/fs"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// 埋め込む静的ファイル群
var (
	//go:embed views/*
	indexFile embed.FS

	//go:embed pb_public/*
	staticFiles embed.FS
)

// ファイル公開
func LoadingStaticFile(pb *pocketbase.PocketBase) error {
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error { // CONTEXT: OnServe.BindFuncで鯖起動時のイベントの関数を設定
		// プロジェクト内のリソースをパスで指定するのではなく、埋め込んだembed.FS型変数で指定

		// indexページを埋め込み
		indexStaticFS, err := fs.Sub(indexFile, "views")
		if err != nil {
			return err
		}
		// CONTEXT: se.Routerはpbが内部的に利用しているHTTPルーターで、ここにエンドポイントを追加できる
		// ルーティング
		se.Router.GET("/{path...}", apis.Static(indexStaticFS, false)) // 第三引数はキャッシュの有無 // "/index/{path...}"

		// 任意のバスに対して適切な静的コンテンツを提供...静的提供ディレクトリの埋め込み
		publicStaticFS, err := fs.Sub(staticFiles, "pb_public")
		if err != nil {
			return err
		}
		// ルーティング
		se.Router.GET("/public/{path...}", apis.Static(publicStaticFS, false))

		return se.Next()
	})

	return nil
}
