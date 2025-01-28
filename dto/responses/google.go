package responses

import (
	"time"
)

// 認証確認
type AuthUserInfo struct {
	Id          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvatarPath  string    `json:"avatorUrl"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	AccessToken string    `json:"accessToken"`
}

// リフレッシュ更新エンドポイント
type RefreshToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}
