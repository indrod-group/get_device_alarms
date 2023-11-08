package auth

type AuthRequest struct {
	Appid     string `json:"appid"`
	Time      int64  `json:"time"`
	Signature string `json:"signature"`
}
