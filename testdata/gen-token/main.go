package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func main() {
	// 1. Lese testdata/private.pem
	data, err := os.ReadFile("testdata/private.pem")
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
	}
	// 2. Bygge en JWT med claims over
	key, err := jwk.ParseKey(data, jwk.WithPEM(true))
	if err != nil {
		log.Fatalf("failed to parse data: %v", err)
	}
	// Les public.pem for å hente kid
	pubData, err := os.ReadFile("testdata/public.pem")
	if err != nil {
		log.Fatalf("failed to public key data: %v", err)
	}
	pubKey, err := jwk.ParseKey(pubData, jwk.WithPEM(true))
	if err != nil {
		log.Fatalf("failed to parse data (public key): %v", err)
	}
	// Generer samme thumbprint som mock-jwks gjør
	thumbprint, err := pubKey.Thumbprint(crypto.SHA256)
	if err != nil {
		log.Fatalf("failed to generate thumbprint: %v", err)
	}
	kid := base64.RawURLEncoding.EncodeToString(thumbprint)
	// Sett kid på privKey
	if err := key.Set(jwk.KeyIDKey, kid); err != nil {
		log.Fatalf("failed to set kid: %v", err)
	}
	if err := key.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("failed to set alg: %v", err)
	}
	t := jwt.New()
	_ = t.Set("iss", "https://gitlab.skead.no")
	_ = t.Set("aud", "min-app-test")
	_ = t.Set("namespace_path", "lagring")
	_ = t.Set("project_path", "lagring/aksjer/sub-aksjer-synt/aksjer-synt-base")
	_ = t.Set("ref", "main")
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(1*time.Hour))
	//_ = t.Set(jwk.KeyIDKey, kid)
	// 3. Signere med RS256
	tokenBytes, err := jwt.Sign(t, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		log.Fatalf("failed to sign jwt:%v", err)
	}
	// 4. Printe tokenet til stdout
	fmt.Println(string(tokenBytes))
}
