package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Config struct {
	JWKSUrl           string
	Audience          string
	AllowedNamespaces []string
}

type Auth struct {
	config Config
}

func NewAuth(config Config) *Auth {
	// returner en ny Auth med config satt
	return &Auth{config: config}
}

func (a *Auth) RequireJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Hent Authorization-header og strip "Bearer " prefikset
		//    r.Header.Get("Authorization") gir deg "Bearer eyJhbG..."
		//    strings.TrimPrefix(..., "Bearer ") gir deg bare tokenet
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, "mangler token", http.StatusUnauthorized)
			return
		}
		// 2. Hent JWKS fra a.config.JWKSUrl
		keySet, err := jwk.Fetch(context.Background(), a.config.JWKSUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 3. Verifiser og parse tokenet
		//    jwt.Parse([]byte(token), jwt.WithKeySet(keySet), jwt.WithAudience(a.config.Audience))
		parsed, err := jwt.Parse([]byte(token), jwt.WithKeySet(keySet), jwt.WithAudience(a.config.Audience))
		//    returnerer (jwt.Token, error) — error betyr ugyldig token
		if err != nil {
			fmt.Printf("parse error: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 4. Hent namespace_path fra token
		npAny, ok := parsed.Get("namespace_path") // returnerer (any, bool)
		if !ok {
			http.Error(w, "namespace_path mangler i token", http.StatusBadRequest)
			return
		}
		np, ok := npAny.(string)
		if !ok {
			http.Error(w, "namespace_path ikke en streng", http.StatusBadRequest)
			return
		}
		// 5. Sjekk at namespace_path er i AllowedNamespaces
		//    loop over a.config.AllowedNamespaces og sammenlign
		for _, ns := range a.config.AllowedNamespaces {
			if strings.HasPrefix(string(np), ns) {
				next(w, r)
				return
			}
		}
		// 6. Kall next hvis alt er ok, ellers http.Error 401
		http.Error(w, "Token does not allow access", http.StatusForbidden)
	}
}
