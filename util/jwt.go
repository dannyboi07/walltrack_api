package util

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"walltrack/common"
	"walltrack/schema"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type jwtCustomClaims struct {
	*jwt.RegisteredClaims
	schema.UserForToken
}

func createToken(userDetails schema.UserForToken, activeDuration time.Duration) (string, int, error) {
	var token *jwt.Token = jwt.New(jwt.GetSigningMethod("RS256"))
	var createdTime time.Time = time.Now()
	var expireAtTime time.Time = createdTime.Add(activeDuration)

	token.Claims = jwtCustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAtTime),
		},
		userDetails,
	}

	signedToken, err := token.SignedString(common.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int(expireAtTime.Sub(createdTime).Seconds()), nil
}

func CreateAccessToken(userDetails schema.UserForToken) (string, int, error) {
	var expireIn time.Duration = time.Minute * 15

	return createToken(userDetails, expireIn)
}

func CreateRefreshToken(userDetails schema.UserForToken) (string, int, error) {
	var expireIn time.Duration = time.Hour * 24 * 7

	return createToken(userDetails, expireIn)
}

func VerifyJwtToken(token string) (jwt.MapClaims, int, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if tokenAlg := t.Method.Alg(); tokenAlg != "RS256" {
			return nil, fmt.Errorf("Unexpected signing method: %s", tokenAlg)
		}

		return common.PublicKey, nil
	})

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if parsedToken.Valid {
		return parsedToken.Claims.(jwt.MapClaims), http.StatusOK, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, http.StatusBadRequest, errors.New("Malformed token")
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, http.StatusUnauthorized, errors.New("token expired")
	}

	Log.Printf("Couldn't handle this JWT token, err: %s, token: %s", err, token)
	return nil, http.StatusInternalServerError, errors.New("Sorry, something went wrong")
}

func ParseJwtClaims(jwtClaims jwt.MapClaims) (map[string]any, int, error) {
	var embeddedDetails map[string]any = make(map[string]any)
	var isMalformed bool = false

	if userIdString, ok := jwtClaims["id"].(string); ok {
		userId, err := primitive.ObjectIDFromHex(userIdString)
		if err != nil {
			isMalformed = true
		}
		embeddedDetails["id"] = userId
	} else {
		isMalformed = true
	}

	if userEmail, ok := jwtClaims["email"].(string); ok {
		embeddedDetails["email"] = userEmail
	} else {
		isMalformed = true
	}

	if isMalformed {
		return nil, http.StatusBadRequest, errors.New("Malformed token")
	}

	return embeddedDetails, 0, nil
}
