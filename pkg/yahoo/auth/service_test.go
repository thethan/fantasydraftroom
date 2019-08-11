package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestYahooAuth_IsValidToken_True(t *testing.T) {
	ya := YahooAuth{
		createdAt: time.Now(),
	}

	nowTime := time.Now()
	ya.UpdateToken("accessToken", "tokenType", "refreshToken", "xoauthGudid", &nowTime, time.Duration(time.Second * 60))

	assert.True(t, ya.IsValidToken(), "could not find a valid time")

	localTime  := time.Local
	nowTime = time.Date(200, 1,1,1,1,1,1, localTime)
	ya.UpdateToken("accessToken", "tokenType", "refreshToken", "xoauthGudid", &nowTime, time.Duration(time.Second * 60))

	assert.False(t, ya.IsValidToken(), "could not find false ")
}