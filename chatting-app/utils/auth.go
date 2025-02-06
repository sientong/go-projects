package utils

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	Role     string `json:"role"`
}

type User struct {
	Id          int64
	Username    string
	Password    string
	Email       string
	FullName    string
	PhoneNumber string
	IsActive    bool
	Role        string
	CreatedAt   string
	UpdatedAt   string
}

func Authenticate(requestData RequestData) (string, error) {

	var authorizedUser User
	query := "SELECT * FROM users WHERE username = $1 and password = $2"
	err := db.QueryRow(query, requestData.Username, requestData.Password).Scan(&authorizedUser.Id, &authorizedUser.Username, &authorizedUser.Password, &authorizedUser.Email, &authorizedUser.FullName, &authorizedUser.PhoneNumber, &authorizedUser.IsActive, &authorizedUser.Role, &authorizedUser.CreatedAt, &authorizedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Invalid username or password")
		} else {
			log.Printf("Failed to query %s", query)
			log.Printf("%s", err)
		}
		return "", err
	}

	tokenString, err := GenerateToken(authorizedUser)
	if err != nil {
		log.Printf("Error on generating token: %s", err)
	}

	return tokenString, err
}

func GenerateToken(authorizedUser User) (string, error) {
	claims := CustomClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    JWT_ISSUER,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Username: authorizedUser.Username,
		Email:    authorizedUser.Email,
		Role:     authorizedUser.Role,
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
