package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

func (config *middlewareConfig) jwtMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, parseToken(config.JWKSet))

			if err != nil || !token.Valid {
				log.Println("error validating token:", err)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			issuer, err := token.Claims.GetIssuer()
			if err != nil || issuer != config.Auth.Url {
				log.Println("error validating issuer:", err)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseToken(jwkSet map[string]*rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		alg := token.Method.Alg()
		publicKey, ok := jwkSet[alg]
		if !ok {
			return nil, fmt.Errorf("no key found for signing method: %v", alg)
		}
		return publicKey, nil
	}
}
