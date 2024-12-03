package main

import (
	"fmt"
	"himaplus-authn/common/logging"
	"himaplus-authn/route"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
)

// 初期化の成果物
type InitInstance struct {
	App *pocketbase.PocketBase
	// Container *dig.Container
}

// mainでの初期化処理
func Init() (*InitInstance, error) {
	// 成果物構造体の宣言
	initInstance := &InitInstance{} // 同じ: initInstance := new(InitInstance)

	// ログ設定を初期化
	err := logging.InitLogging() // セットアップ
	if err != nil {              // エラーチェック
		fmt.Printf("error set up logging: %v\n", err) // ログ関連のエラーなのでログは出力しない
		panic("error set up logging.")
	}
	logging.SuccessLog("Start server!")

	// .envから定数をプロセスの環境変数にロード
	err = godotenv.Load(".env") // エラーを格納
	if err != nil {             // エラーがあったら
		logging.ErrorLog("Error loading .env file.", err)
		return nil, err
	}

	// DB初期化やルーティング設定など、依存関係にかかわるものの初期化とDIコンテナによる各層の依存関係登録
	initInstance.App = route.Routing()

	return initInstance, nil
}
