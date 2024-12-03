package main

import (
	"himaplus-authn/common/logging"
)

func main() {
	// 初期化処理
	initInstances, err := Init() // add "initInstances, " when changing to ddd
	if err != nil {
		return
	}
	// 破棄処理
	defer logging.LogFile().Close() // defer文でこの関数終了時に破棄
	logging.SuccessLog("Successful server init process.")

	// 鯖起動
	if err := initInstances.App.Start(); err != nil {
		logging.ErrorLog("Failed to start server", err)
		panic(err)

	}
}
