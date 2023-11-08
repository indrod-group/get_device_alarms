package auth

type AuthResponse struct {
	Code        int     `json:"code"`
	Expiresin   *int    `json:"expiresIn,omitempty"`
	AccessToken *string `json:"accessToken,omitempty"`
	Result      *string `json:"result,omitempty"`
}
