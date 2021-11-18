package types

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

type CanisterInfo struct {
	Repo            string `json"repo"`
	Owner           string `json:"owner"`
	Controller      string `json:"controller"`
	CanisterName    string `json:"name"`
	CanisterID      string `json:"id"`
	CanisterType    string `json:"type"`
	Framework       string `json:"framework"`
	Network         string `json:"network"`
	CreateTimestamp int64  `json:"createtimestamp"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
}
