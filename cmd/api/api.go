package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/yofabr/mono-client/cmd/application"
	"github.com/yofabr/mono-client/internal/auth"
	"github.com/yofabr/mono-client/internal/middleware"
)

// Api wires HTTP routes to the underlying application services.
type Api struct {
	app *application.Application
}

// NewApi creates the HTTP API adapter using the shared application container.
func NewApi(app application.Application) *Api {
	api := Api{app: &app}
	return &api
}

// Init registers all HTTP handlers on the default net/http mux.
func (api *Api) Init() {
	// Health/root route that returns the resolved client IP.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		msg := "Your IP Address is: " + clientIP
		_, err := w.Write([]byte(msg + "\n"))
		if err != nil {
			return
		}
	})

	authHandler := auth.NewAuthHandler(api.app)
	// Allow max 5 auth attempts per minute for each route+IP key.
	authRateLimiter := middleware.NewRateLimiter(5, time.Minute)

	// Login endpoint: validates payload, then issues token on success.
	http.HandleFunc("/login", authRateLimiter.Handler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var creds auth.Credentials
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // Reject unexpected JSON keys early.

		if dec.Decode(&creds) != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ip := getClientIP(r)
		res, err := authHandler.Login(ip, creds.Email, creds.Password)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(res + "\n"))
		if err != nil {
			return
		}
	}, func(r *http.Request) string {
		return "/login:" + getClientIP(r)
	}))

	// Register endpoint: creates user, stores active auth session, returns token.
	http.HandleFunc("/register", authRateLimiter.Handler(func(w http.ResponseWriter, r *http.Request) {
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
		res, err := authHandler.Register(ip, creds.Email, creds.Password)
		if err != nil {
			return
		}

		_, err = w.Write([]byte(res + "\n"))
		if err != nil {
			return
		}
	}, func(r *http.Request) string {
		return "/register:" + getClientIP(r)
	}))
}

// getClientIP extracts the original caller IP from proxy headers or socket addr.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
