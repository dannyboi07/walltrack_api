package util

import (
	"net/http"
	"walltrack/common"
	"walltrack/schema"
)

func CreateAccessTokenCookie(userDetails schema.UserForToken) (*http.Cookie, error) {
	var (
		accessToken string
		expiresIn   int
		err         error
	)
	accessToken, expiresIn, err = CreateAccessToken(userDetails)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     common.AccessTokenCookieName,
		Value:    accessToken,
		MaxAge:   expiresIn,
		Path:     common.AccessTokenCookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}, nil
}

func CreateRefreshTokenCookie(userDetails schema.UserForToken) (*http.Cookie, error) {
	var (
		refreshToken string
		expiresIn    int
		err          error
	)
	refreshToken, expiresIn, err = CreateRefreshToken(userDetails)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     common.RefreshTokenCookieName,
		Value:    refreshToken,
		MaxAge:   expiresIn,
		Path:     common.RefreshTokenCookiePath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}, nil
}
