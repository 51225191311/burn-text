package internal

type CreateSecretReq struct {
	Text     string `json:"text" binding:"required"`
	Password string `json:"password"` //可选密码
}

type SecretData struct {
	CipherText   string `json:"cipher_text"`
	PasswordHash string `json:"password_hash,omitempty"` //可选密码,否则为空
}

type CreateSecretResp struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}
