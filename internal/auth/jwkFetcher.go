package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
)

func (config *middlewareConfig) getJWKSet() error {
	// Make the GET request
	response, err := http.Get(config.Auth.Url + "/protocol/openid-connect/certs")
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body:", err)
		}
	}(response.Body)

	// Decode the JSON response
	var jwkSet struct {
		Keys []struct {
			Kid string   `json:"kid"`
			Alg string   `json:"alg"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5C []string `json:"x5c"`
		} `json:"keys"`
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jwkSet); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	// Create a map to store RSA public keys
	jwkMap := make(map[string]*rsa.PublicKey)

	// Iterate through each key in the JWK set
	for _, key := range jwkSet.Keys {
		// Decode base64url-encoded modulus (N) and exponent (E)
		modulus, err := decodeBase64URL(key.N)
		if err != nil {
			return fmt.Errorf("error decoding modulus: %v", err)
		}

		exponent, err := decodeBase64URL(key.E)
		if err != nil {
			return fmt.Errorf("error decoding exponent: %v", err)
		}

		// Create RSA public key
		pubKey := &rsa.PublicKey{
			N: modulus,
			E: int(exponent.Int64()),
		}

		// Store the public key in the map using the key ID (Kid)
		jwkMap[key.Alg] = pubKey
	}

	config.JWKSet = jwkMap

	return nil
}

func decodeBase64URL(input string) (*big.Int, error) {
	base64Str := strings.ReplaceAll(input, "-", "+")
	base64Str = strings.ReplaceAll(base64Str, "_", "/")

	switch len(base64Str) % 4 {
	case 2:
		base64Str += "=="
	case 3:
		base64Str += "="
	}

	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	result := new(big.Int).SetBytes(data)
	return result, nil
}
