package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

type RequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Group    string `json:"group"`
}

func Authenticate(requestData RequestData) (bool, M) {
	basePath, _ := os.Getwd()
	dbPath := filepath.Join(basePath, "users.json")
	buf, _ := os.ReadFile(dbPath)

	data := make([]M, 0)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return false, nil
	}

	res := gubrak.From(data).Find(func(each M) bool {
		return each["username"] == requestData.Username && each["password"] == requestData.Password
	}).Result()

	if res != nil {
		resM := res.(M)
		delete(resM, "password")
		return true, resM
	}

	return false, nil
}

func GenerateToken(userInfo M) (error, []byte) {
	claims := CustomClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    JWT_ISSUER,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Username: userInfo["username"].(string),
		Email:    userInfo["email"].(string),
		Group:    userInfo["group"].(string),
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return err, nil
	}

	tokenString, _ := json.Marshal(M{"token": signedToken})
	return nil, tokenString
}
