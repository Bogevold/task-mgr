package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bogevold/task-mgr/internal/auth"
	"github.com/bogevold/task-mgr/internal/handler"
	pgstore "github.com/bogevold/task-mgr/internal/store"
	"github.com/bogevold/task-mgr/internal/task"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	storeType := os.Getenv("STORE")
	var store task.TaskStore
	var err error
	switch storeType {
	case "memory":
		fmt.Println("Startet i debug modus med in memory database")
		store = task.NewInMemoryStore()
		for i := range 4 {
			_, err := store.Save(task.Task{Title: fmt.Sprintf("Task %d", i), CreatedAt: time.Now()})
			if err != nil {
				fmt.Printf("forventet ingen feil, fikk: %v", err)
			}
		}
	default:
		fmt.Println("Startet i standard modus med in PostgreSQL database")
		connStr := os.Getenv("DATABASE_URL")
		store, err = pgstore.NewPostgresStore(connStr)
		if err != nil {
			fmt.Printf("Kunne ikke initialisere postgres backend: %v", err)
			return
		} // bruk PostgresStore

	}

	// Lese JWKS_URL, JWT_AUDIENCE og ALLOWED_NAMESPACES fra miljøvariabler
	jwksUrl := os.Getenv("JWKS_URL")
	jwtAud := os.Getenv("JWT_AUDIENCE")
	alwdNs := strings.Split(os.Getenv("ALLOWED_NAMESPACES"), ",")
	// Opprette en auth.Config
	// Opprette en auth.Auth med auth.NewAuth
	newAuth := auth.NewAuth(auth.Config{JWKSUrl: jwksUrl, Audience: jwtAud, AllowedNamespaces: alwdNs})
	// Sende den inn i RegisterRoutes

	th := handler.NewTaskHandler(store)
	mux := http.NewServeMux()
	th.RegisterRoutes(mux, newAuth)
	port := getEnv("PORT", "8072")
	fmt.Printf("Server lytter på :%s\n", port)
	//http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		fmt.Printf("Server feilet: %v\n", err)
	}
}
