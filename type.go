package main

type AuthResponse struct {
	AccessToken          string `json:"access_token"`
	ExpiresIn            int16  `json:"expires_in"`
	RefreshTOken         string `json:"refresh_token"`
	RefreshTokenExpireIn int16  `json:"refresh_token_expires_in"`
	Scope                string `json:"scope"`
	TokenType            string `json:"token_type"`
}

type CmdOutput struct {
	TaskID string
	Logs   []string
}
