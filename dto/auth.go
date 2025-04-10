package dto

type AuthRequest struct {
	AppId        string `json:"appId"`
	ClientSecret string `json:"clientSecret"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}
