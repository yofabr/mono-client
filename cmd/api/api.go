package api

import (
	"encoding/json"
	"errors"
	"io"
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
	// Health liveness probe - app is running
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	// Readiness probe - app can handle requests (DB and Redis connected)
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := api.app.Databases.PostgresHealth(ctx); err != nil {
			http.Error(w, "postgres not ready", http.StatusServiceUnavailable)
			return
		}

		if err := api.app.Databases.RedisHealth(ctx); err != nil {
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

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

		creds, err := decodeCredentials(w, r)
		if err != nil {
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

		creds, err := decodeCredentials(w, r)
		if err != nil {
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

func decodeCredentials(w http.ResponseWriter, r *http.Request) (auth.Credentials, error) {
	var creds auth.Credentials
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&creds); err != nil {
		return auth.Credentials{}, err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return auth.Credentials{}, errors.New("unexpected trailing data")
	}

	return creds, nil
}

// getClientIP extracts the original caller IP from proxy headers or socket addr.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		if net.ParseIP(xrip) != nil {
			return xrip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if net.ParseIP(r.RemoteAddr) != nil {
			return r.RemoteAddr
		}
		return ""
	}

	if net.ParseIP(host) == nil {
		return ""
	}

	return host
}
