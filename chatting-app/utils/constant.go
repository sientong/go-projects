package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SERVER_HOST = "localhost"
var SERVER_PORT = "8080"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNATURE_KEY = []byte("secret")
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_ISSUER = "chatting-app"
var EXCLUDED_PATH = []string{"/login", "/ws", "/chat", "/"}
var MESSAGE_NEW_USER = "New User"
var MESSAGE_CHAT = "Chat"
var MESSAGE_LEAVE = "Leave"
var BASE_FILE_PATH = "views"
