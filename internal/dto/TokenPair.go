package dto

type TokenPair struct {
	RefreshToken string `json:"refresh_token" example:"refresh_base64_token"`
	AccessToken string `json:"access_token" example:"jwt_access_token"`
}