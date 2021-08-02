package auth

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// WebexAccessTokenResponse is the response from Webex on requesting an access token using authorization code grant
type WebexAccessTokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Message               string
	Errors                []struct {
		Description string `json:"description"`
	} `json:"errors"`
	TrackingID string `json:"trackingId"`
}

func (t WebexAccessTokenResponse) String() string {
	return fmt.Sprintf(`Token: %s
ExpiresIn: %d
RefreshToken: %s
RefreshTokenExpiresIn: %d
Message: %s
Errors: %s
TrackingID: %s`, t.AccessToken, t.ExpiresIn, t.RefreshToken, t.RefreshTokenExpiresIn, t.Message, t.Errors, t.TrackingID)
}

// Save saves the token details to the config file
func (t WebexAccessTokenResponse) Save() error {
	expiration := time.Now().Unix() + int64(t.ExpiresIn)
	refreshExpiration := time.Now().Unix() + int64(t.RefreshTokenExpiresIn)
	viper.Set("expiration", expiration)
	viper.Set("refresh_expiration", refreshExpiration)
	viper.Set("token", t.AccessToken)
	viper.Set("refresh_token", t.RefreshToken)
	return viper.WriteConfig()
}

// Expires returns a string representation of when the token actually expires.  Bear in mind
// that it is based off the ExpiresIn parameter which is relative to when you received the token.
func (t WebexAccessTokenResponse) Expires() string {
	expires := time.Now().Unix() + int64(t.ExpiresIn)
	return time.Unix(expires, 0).Format(time.RFC3339)
}
