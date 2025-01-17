package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type CustomMux struct {
	http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func (c *CustomMux) RegisterMiddleware(middleware func(http.Handler) http.Handler) {
	c.middlewares = append(c.middlewares, middleware)
}

func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &c.ServeMux

	for _, middleware := range c.middlewares {
		current = middleware(current)
	}

	current.ServeHTTP(w, r)
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("url " + r.URL.RequestURI() + " path " + r.URL.Path + " method " + r.Method)
		if len(EXCLUDED_PATH) > 0 && contains(EXCLUDED_PATH, r.URL.Path) {
			log.Println("Path is excluded from JWT validation")
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Please Login First", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("signing method invalid")
			}

			return JWT_SIGNATURE_KEY, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}

		type contextKey string
		var userContextKey contextKey = "userInfo"
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func contains(arr []string, str string) bool {
	for _, each := range arr {
		if each == str {
			return true
		}
	}

	return false
}
