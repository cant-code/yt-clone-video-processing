package auth

import (
	"log"
	"net/http"
)

func HandleJwtAuthMiddleware() func(http.Handler) http.Handler {
	set, err := getJWKSet("http://localhost:8900/realms/yt-clone/protocol/openid-connect/certs")
	if err != nil {
		log.Printf("Error fetching jwk-sets: %v\n", err)
	}

	return jwtMiddleware(set)
}
