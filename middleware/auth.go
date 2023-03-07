package middleware

import (
	"context"
	"net/http"
	"walltrack/util"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			accessTokenCookie *http.Cookie
			err               error
		)
		accessTokenCookie, err = r.Cookie("accessToken")
		if err != nil {
			util.WriteApiErrMessage(w, http.StatusUnauthorized, "Missing access token")
			return
		}

		var (
			jwtClaims  jwt.MapClaims
			statusCode int
		)
		jwtClaims, statusCode, err = util.VerifyJwtToken(accessTokenCookie.Value)
		if err != nil {
			util.WriteApiErrMessage(w, statusCode, err.Error())
			return
		}

		var userDetails map[string]interface{}
		userDetails, statusCode, err = util.ParseJwtClaims(jwtClaims)
		if err != nil {
			util.WriteApiErrMessage(w, statusCode, err.Error())
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userDetails", userDetails)))
	})
}
