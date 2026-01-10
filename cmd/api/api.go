package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/yofabr/mono-client/cmd/application"
	"github.com/yofabr/mono-client/internal/auth"
)

type Api struct {
	app *application.Application
}

func NewApi(app application.Application) *Api {
	api := Api{
		app: &app,
	}
	return &api
}

func (api *Api) Init() {
	// Root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		msg := "Your IP is: " + clientIP
		w.Write([]byte(msg + "\n"))
	})

	authHandler := auth.NewAuthHandler(api.app)

	// Login handler
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds auth.Credentials
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if dec.Decode(&creds) != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		ip := getClientIP(r)
		authHandler.Login(ip, creds.Email, creds.Password)
		log.Println("Login attempt:", creds.Email, creds.Password)
		w.Write([]byte(fmt.Sprintf("Login received for: %s", creds.Email)))
	})

	// Register handler
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds auth.Credentials
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if dec.Decode(&creds) != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ip := getClientIP(r)
		authHandler.Register(ip, creds.Email, creds.Password)

		log.Println("Register attempt:", creds.Email, creds.Password)
		w.Write([]byte(fmt.Sprintf("Register received for: %s", creds.Email)))
	})
	// log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}

func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For first (if behind a proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
