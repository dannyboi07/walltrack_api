package controller

import (
	"encoding/json"
	"net/http"
	"walltrack/common"
	"walltrack/db"
	"walltrack/model"
	"walltrack/schema"
	"walltrack/util"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {
	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userRegister schema.UserRegister
		statusCode   int
		err          error
	)
	statusCode, err = util.JsonParseErr(jDec.Decode(&userRegister))
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Failed to decode req, err:", err)
		return
	}

	statusCode, err = userRegister.Validate()
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Register user request validation err:", err)
		return
	}

	var hashedPw string
	hashedPw, err = util.HashPassword(*userRegister.Password)
	if err != nil {
		util.WriteApiErrMessage(w, 0, "")
		util.Log.Println("Failed to hash password, err:", err)
		return
	}

	_, err = db.InsertUser(model.User{
		Email:        *userRegister.Email,
		PasswordHash: hashedPw,
	})
	if err != nil {
		statusCode, message := 0, ""
		if db.IsDuplicateKeyError(err) {
			statusCode, message = http.StatusBadRequest, "Account already exists"
			util.Log.Println("Registration attempt with existing email:", *userRegister.Email, r.RemoteAddr)
		} else {
			util.Log.Println("Failed to register/insert user details into db, err:", err)
		}
		util.WriteApiErrMessage(w, statusCode, message)
		return
	}

	util.WriteApiMessage(w, 0, "Your account has been created!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	jDec := json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userLogin  schema.UserLogin
		statusCode int
		err        error
	)
	statusCode, err = util.JsonParseErr(jDec.Decode(&userLogin))
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Failed to decode req, err:", err)
		return
	}

	statusCode, err = userLogin.Validate()
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("User login request validation err:", err)
		return
	}

	user, found, err := db.SelectUserByEmail(*userLogin.Email)
	if err != nil {
		util.WriteApiErrMessage(w, 0, "")
		util.Log.Println("Failed to get user details for login, err:", err)
		return
	} else if !found {
		util.WriteApiErrMessage(w, http.StatusUnauthorized, "Check your email/password & try again")
		return
	}

	if isCorrectPw := util.VerifyPassword(user.PasswordHash, *userLogin.Password); !isCorrectPw {
		util.WriteApiErrMessage(w, http.StatusUnauthorized, "Check your email/password & try again")
		return
	}

	var userTokenDetails schema.UserForToken = schema.UserForToken{
		Id:    user.Id,
		Email: user.Email,
	}

	var accessTokenCookie *http.Cookie
	accessTokenCookie, err = util.CreateAccessTokenCookie(userTokenDetails)
	if err != nil {
		util.WriteApiErrMessage(w, 0, "")
		util.Log.Println("Failed to create access token")
		return
	}
	var refreshTokenCookie *http.Cookie
	refreshTokenCookie, err = util.CreateRefreshTokenCookie(userTokenDetails)
	if err != nil {
		util.WriteApiErrMessage(w, 0, "")
		util.Log.Println("Failed to create refresh token")
		return
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	util.WriteApiMessage(w, 0, "Logged in")
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var (
		refreshTokenCookie *http.Cookie
		err                error
	)
	refreshTokenCookie, err = r.Cookie(common.RefreshTokenCookieName)
	if err != nil {
		util.WriteApiErrMessage(w, http.StatusUnauthorized, "Session expired")
		return
	}

	var (
		jwtClaims  jwt.MapClaims
		statusCode int
	)
	jwtClaims, statusCode, err = util.VerifyJwtToken(refreshTokenCookie.Value)
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Failed to verify refresh token, err:", err, r.RemoteAddr)
		return
	}

	var userDetails map[string]interface{}
	userDetails, statusCode, err = util.ParseJwtClaims(jwtClaims)
	if err != nil {
		util.WriteApiErrMessage(w, statusCode, err.Error())
		util.Log.Println("Failed to parse JWT claims, err:", err, r.RemoteAddr)
		return
	}

	var userTokenDetails schema.UserForToken = schema.UserForToken{
		Id:    userDetails["id"].(primitive.ObjectID),
		Email: userDetails["email"].(string),
	}
	var accessTokenCookie *http.Cookie
	accessTokenCookie, err = util.CreateAccessTokenCookie(userTokenDetails)
	if err != nil {
		util.WriteApiErrMessage(w, 0, "")
		util.Log.Println("Failed to create access token")
		return
	}

	http.SetCookie(w, accessTokenCookie)
	w.WriteHeader(http.StatusOK)
}
