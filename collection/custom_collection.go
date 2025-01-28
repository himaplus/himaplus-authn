package collection

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types" // "github.com/pocketbase/pocketbase/tools/types"
)

// コレクションの操作
func CustomCollection(pb *pocketbase.PocketBase) *pocketbase.PocketBase {
	// access token列を追加
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// カスタム対象のコレクションを取得
		userCollection, err := se.App.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		// change rule types.Pointer(""): 誰でも, types.Pointer("@request.auth.id != ''"): 認証済みユーザーのみ, types.Pointer("@request.auth.id != id"): 認証済みユーザー自身のみ, nil: 管理者のみ
		userCollection.ListRule = nil
		userCollection.ViewRule = types.Pointer("@request.auth.id != id")
		userCollection.CreateRule = types.Pointer("@request.auth.id != id") // 確定
		userCollection.UpdateRule = types.Pointer("@request.auth.id != id")
		userCollection.DeleteRule = nil

		// フィールドやリレーションキーを追加

		// リフレッシュトークンのテキストフィールド追加
		userCollection.Fields.Add(&core.TextField{
			Name:     "refreshToken",
			Max:      255, // The number of characters in google access token is 222.
			Required: false,
		})

		// アクセストークンのテキストフィールド追加　HACK: いらんかも
		userCollection.Fields.Add(&core.TextField{
			Name:     "accessToken",
			Max:      255, // The number of characters in google refresh token is 103.
			Required: false,
		})

		// 保存
		err = se.App.Save(userCollection)
		if err != nil {
			return err
		}

		return se.Next()
	})

	return pb
}
