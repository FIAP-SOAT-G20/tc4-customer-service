package presenter

import "encoding/json"

type JWTResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int64  `json:"expires_in" example:"3600"`
}

func (r JWTResponse) String() string {
	o, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(o)
}
