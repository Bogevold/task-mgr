package main

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

func main() {
	// 1. Les public.pem
	data, err := os.ReadFile("testdata/public.pem")
	if err != nil {
		log.Fatalf("failed to read public key: %v", err)
	}
	// 2. Konverter til JWKS
	key, err := jwk.ParseKey(data, jwk.WithPEM(true))
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	thumbprint, err := key.Thumbprint(crypto.SHA256)
	if err != nil {
		log.Fatalf("failed to generate tumpringt: %v", err)
	}
	if err := key.Set(jwk.KeyIDKey, base64.RawURLEncoding.EncodeToString(thumbprint)); err != nil {
		log.Fatalf("failed to set kid: %v", err)
	}
	if err := key.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("failed to set alg: %v", err)
	}
	// if err := key.Set(jwk.KeyIDKey, "my-key"); err != nil {
	//   log.Fatalf("failed to set kid: %v", err)
	// }
	set := jwk.NewSet()
	set.AddKey(key)

	http.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(set); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// 3. Start server på :8090
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
