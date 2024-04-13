package auth

import (
	"crypto/rsa"
	"log"
	"net/http"
	"yt-clone-video-processing/internal/configurations"
)

type IMiddleware interface {
	getJWKSet() error
	jwtMiddleware() func(http.Handler) http.Handler
}

type middlewareConfig struct {
	Auth   configurations.Auth
	JWKSet map[string]*rsa.PublicKey
}

func HandleJwtAuthMiddleware(auth *configurations.Auth) func(http.Handler) http.Handler {
	middleware := IMiddleware(&middlewareConfig{Auth: *auth})

	err := middleware.getJWKSet()
	if err != nil {
		log.Printf("Error fetching jwk-sets: %v\n", err)
	}

	return middleware.jwtMiddleware()
}
