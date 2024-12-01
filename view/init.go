package view

import (
	"embed"
	"io/fs"
)

// 埋め込む静的ファイル群
//
//go:embed views/*
var indexFile embed.FS

//go:embed pb_public/*
var staticFiles embed.FS

// indexページ
func EmbedIndexFile() (fs.FS, error) {
	return fs.Sub(indexFile, "views") // プロジェクト内のリソースをパスで指定するのではなく、埋め込んだembed.FS型変数で指定
}

// 静的提供ファイル
func EmbedStaticFile() (fs.FS, error) {
	return fs.Sub(staticFiles, "pb_public")
}
