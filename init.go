package main

import (
	"fmt"
	"himaplus-authn/collection"
	"himaplus-authn/common/logging"
	"himaplus-authn/hook"
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
	initInstance := &InitInstance{
		App: pocketbase.New(), // pbインスタンス
	} // initInstance := new(InitInstance)も同じだが、リテラル構文だとフィールドの初期化ができる

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

	// initInstanceを初期化
	initInstance.App = collection.CustomCollection(initInstance.App) // コレクションのカスタム
	initInstance.App = hook.RegisterHooks(initInstance.App) // 既存APIに対してのフック処理
	initInstance.App = route.SetupRouter(initInstance.App) // ルーティング設定など

	return initInstance, nil
}
