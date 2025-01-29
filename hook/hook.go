package hook

// 既存APIに対してのフック処理

import (
	"fmt"
	"himaplus-authn/common/logging"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterHooks(pb *pocketbase.PocketBase) *pocketbase.PocketBase {
	// リフレッシュの検証
	pb.OnRecordAuthRefreshRequest("users").BindFunc(func(re *core.RecordAuthRefreshRequestEvent) error {
		// 更新
		logging.SimpleLog("re.Collection.AuthToken: ", re.Collection.AuthToken)
		return re.Next()
	})

	// 認証リクエスト後にrefresh tokenを保存
	pb.OnRecordAuthWithOAuth2Request("users").BindFunc(func(re *core.RecordAuthWithOAuth2RequestEvent) error {
		// before処理

		// log
		logging.SimpleLog("Here is OnRecordAuthWithOAuth2Request before hook.")

		if err := re.Next(); err != nil { // before end: return re.Next()
			logging.ErrorLog("OnRecordAuthWithOAuth2Request: re.Next()", err)
			return err
		}

		// after処理

		// log
		logging.SimpleLog("re.OAuth2User.AccessToken: ", re.OAuth2User.AccessToken)
		logging.SimpleLog("re.OAuth2User.RefreshToken: ", re.OAuth2User.RefreshToken) // access_type=offlineを指定してないとonlineとみなされてもらえない
		// logging.SimpleLog("re: ", re) // ok
		// logging.SimpleLog("re.Auth.Id: ", re.Auth.Id) // ng
		// logging.SimpleLog("re.Collection.Id: ", re.Collection.Id) // ok

		// トークンをカスタムコレクションで追加した列に保存する

		// レコードを取得
		logging.SimpleLog("re.Record.Id: ", re.Record.Id) // ok
		record, err := pb.FindRecordById("users", re.Record.Id)
		if err != nil {
			return err
		}
		fmt.Printf("record: %v\n", record)

		// トークンをレコードに追加
		record.Set("refreshToken", re.OAuth2User.RefreshToken)
		record.Set("accessToken", re.OAuth2User.AccessToken)

		// 保存
		err = pb.Save(record)
		if err != nil {
			return err
		}
		fmt.Printf("record: %v\n", record)

		return nil
	})
	return pb
}
