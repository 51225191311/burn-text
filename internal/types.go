package internal

type CreateSecretReq struct {
	Text string `json:"text" binding:"required"`
}

type CreateSecretResp struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}
