package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"himaplus-authn/common/logging"
	"himaplus-authn/dto/responses"
	"himaplus-authn/view"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types" // "github.com/pocketbase/pocketbase/tools/types"
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

// コレクションの操作
func collectionOpe(pb *pocketbase.PocketBase) {
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
}

// カスタムエンドポイント
func endpointRouting(pb *pocketbase.PocketBase) {
	// 認証状態を確認
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/auth/user", func(re *core.RequestEvent) error {
			// 認証済みユーザー（:アクセスを許されたユーザー）のレコードを取得authUserRecord
			userRecord := re.Auth
			fmt.Printf("authUserRecord: %v\n", userRecord)

			// 構造体にマッピング
			auiResp := &responses.AuthUserInfo{
				Id:    userRecord.Id,
				Email: userRecord.Email(),
				Name:  userRecord.Collection().Name,
			}
			fmt.Printf("auiResp: %v\n", auiResp)

			// 値をJSON形式で返却
			return re.JSON(http.StatusOK, map[string]any{
				"data": map[string]any{
					"authenticated": true,
					"userInfo":      auiResp,
				},
				"message": "The request requires valid record authorization token.",
				"status":  200,
			})
		}).Bind(apis.RequireAuth("_superusers", "users")) // HTTPメソッド関数にチェーンしてミドルウェアを追加できる
		return se.Next()
	})

	// アクセストークンを返す
	pb.OnServe().BindFunc(func(se *core.ServeEvent) error { // pbインスタンスのOnServe()フックを使って処理を鯖起動時にトリガーする // 処理はBindFunc()に渡す
		se.Router.GET("/google/access_token", func(re *core.RequestEvent) error { // core.RequestEventはreq, resを操作するためのメソッドを持つ
			// 認証済みユーザー（:アクセスを許されたユーザー）のレコードを取得authUserRecord
			userRecord := re.Auth
			fmt.Printf("authUserRecord: %v\n", userRecord)

			// 認証済みユーザーのリフレッシュトークンを取得
			refreshToken := userRecord.GetString("refreshToken") // refreshToken := userRecord.FieldsData()["refreshToken"] // 型アサーションが必要
			fmt.Printf("accessToken: %v\n", refreshToken)

			// クライアントIDとクライアントシークレットの取得

			// client_secret_*.jsonを取得
			pattern := "./client_secret_*.json"
			clientSecretJsonFileNames, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			fmt.Printf("clientSecretJsonFileNames: %v\n", clientSecretJsonFileNames)
			slices.Sort(clientSecretJsonFileNames)
			clientSecretJsonFile, err := ioutil.ReadFile(clientSecretJsonFileNames[0])
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			// mapに変換
			var clientSecretMap map[string]any // 空のJSON宣言
			err = json.Unmarshal(clientSecretJsonFile, &clientSecretMap)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			fmt.Printf("clientSecretMap: %v\n", clientSecretMap)
			clientSecretMapWeb, ok := clientSecretMap["web"].(map[string]any)
			if !ok {
				err = errors.New("assertion of value from client_secret.json file failed: web")
				fmt.Printf("err: %v\n", err)
				return err
			}
			// anyをstringアサーションして取り出す
			clientId, clientIdOk := clientSecretMapWeb["client_id"].(string)
			clientSecret, clientSecretOk := clientSecretMapWeb["client_secret"].(string)
			if !clientIdOk || !clientSecretOk {
				err = errors.New("assertion of value from client_secret.json file failed: client_id || client_secret")
				fmt.Printf("err: %v\n", err)
				return err
			}
			fmt.Printf("clientId: %v\n", clientId)
			fmt.Printf("clientSecret: %v\n", clientSecret)

			// Googleへ更新リクエストを送る

			// リクエストの作成
			method := "POST"                                   // メソッド
			endopoint := "https://oauth2.googleapis.com/token" // URL
			form := url.Values{}                               // フォームデータやHTTPクエリパラメータを扱うmap[string][]string型
			form.Set("grant_type", "refresh_token")            // SetはAddと違って同じキーを上書きする
			form.Set("refresh_token", refreshToken)
			form.Set("client_id", clientId)
			form.Set("client_secret", clientSecret)
			body := strings.NewReader(form.Encode())              // HTTPエンコードしてボディを作る
			requ, err := http.NewRequest(method, endopoint, body) // リクエストの作成
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			requ.Header.Set("Content-Type", "application/x-www-form-urlencoded") // ヘッダーを追加

			// リクエストを送る
			client := &http.Client{ // クライアントを作成
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(requ) // リクエストを送信しレスポンスを受け取る
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			defer resp.Body.Close() // リソースの解放 必須

			// 構造体にマッピング
			var refTokenResp responses.RefreshToken
			err = json.NewDecoder(resp.Body).Decode(&refTokenResp)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return err
			}
			fmt.Printf("req.AccessToken: %v\n", refTokenResp.AccessToken)

			// 値をJSON形式で返却
			return re.JSON(http.StatusOK, map[string]string{
				"googleAccessToken": refTokenResp.AccessToken,
			})
		}).Bind(apis.RequireAuth("_superusers", "users")) // HTTPメソッド関数にチェーンしてミドルウェアを追加できる

		return se.Next()
	})

	// リフレッシュの検証
	// fires only for "users" and "managers" auth collections
	pb.OnRecordAuthRefreshRequest("users").BindFunc(func(re *core.RecordAuthRefreshRequestEvent) error {
		// 更新
		fmt.Printf("re.Collection.AuthToken: %v\n", re.Collection.AuthToken) // これを使って
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
}

// ルーティング
func Routing() *pocketbase.PocketBase {
	// pbインスタンス
	pb := pocketbase.New()

	// 拡張

	globalMiddleware(pb) // ミドルウェア
	fileRouting(pb)      // ファイル公開
	collectionOpe(pb)    // コレクション操作
	endpointRouting(pb)  // カスタムエンドポイントのルーティング

	return pb
}
