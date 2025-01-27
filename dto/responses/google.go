package responses

// 認証確認
type AuthUserInfo struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	// AvatorUrl   string `json:"avatorUrl"`
	// Created     int    `json:"created"` // TODO: 型
	// Updated     int    `json:"updated"` // TODO: 型
	// AccessToken string `json:"accessToken"`
}

// リフレッシュ更新エンドポイント
type RefreshToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}
