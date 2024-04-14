package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"yt-clone-video-processing/internal/configurations"
)

type IMiddleware interface {
	getJWKSet() error
	jwtMiddleware() func(http.Handler) http.Handler
}

type openIdConfig struct {
	Issuer string `json:"issuer"`
	Jwks   string `json:"jwks_uri"`
}

type middlewareConfig struct {
	OpenIdConfig *openIdConfig
	JWKSet       map[string]*rsa.PublicKey
}

const wellKnownConfigs = "/.well-known/openid-configuration"

func HandleJwtAuthMiddleware(auth *configurations.Auth) func(http.Handler) http.Handler {
	openIdConfig, err := getOpenIdConfigs(auth)
	if err != nil {
		log.Println("Error getting openid configs: ", err)
	}

	middleware := IMiddleware(&middlewareConfig{OpenIdConfig: openIdConfig})

	err = middleware.getJWKSet()
	if err != nil {
		log.Printf("Error fetching jwk-sets: %v\n", err)
	}

	return middleware.jwtMiddleware()
}

func getOpenIdConfigs(auth *configurations.Auth) (*openIdConfig, error) {
	response, err := http.Get(auth.Url + wellKnownConfigs)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(response.Body)

	var openIdConfig openIdConfig
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&openIdConfig); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &openIdConfig, nil
}
