package auth

import "time"

type YahooAuth struct {
	AccessToken     string
	TokenType       string
	ExpiresIn       time.Duration
	RefreshToken    string
	XoauthYahooGuid string
	createdAt       time.Time
	updatedAt       *time.Time
}

// IsValidToken tells you if the token is valid
func (y *YahooAuth) IsValidToken() bool {
	tokenTimeCreated := y.createdAt
	if y.updatedAt != nil {
		tokenTimeCreated = *y.updatedAt
	}
	currentTime := time.Now()
	tokenTime := currentTime.Sub(tokenTimeCreated)
	return tokenTime < y.ExpiresIn
}

// UpdateToken will update a yahoo token
func (y *YahooAuth) UpdateToken(accessToken, tokenType, refreshToken, XoauthYahooGuid string, updatedAt *time.Time, expiresIn time.Duration) {
	y.AccessToken = accessToken
	y.TokenType = tokenType
	y.RefreshToken = refreshToken
	y.XoauthYahooGuid = XoauthYahooGuid
	y.updatedAt = updatedAt
	y.ExpiresIn = expiresIn

}
