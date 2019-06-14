package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	// gorillaContex "github.com/gorilla/context"
	// "github.com/gorilla/sessions"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"

	_ "github.com/lib/pq"

	"github.com/color-book/web_server/api"
	"github.com/color-book/web_server/configVars"
	"github.com/color-book/web_server/dataStore"
	"github.com/color-book/web_server/handlers"
	"github.com/color-book/web_server/sessionStore"
)

var (
	listenAddr string
)

type MiddlewareException struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Redirect     string `json:"redirect"`
}

func runServer() {
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")
	flag.Parse()

	// INITIALIZE CONFIG VARIABLES
	configVars.InitConfigVars()

	// LOGGING
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Colorbook server is starting...")

	// INITIALIZE DATABASE
	connString := configVars.Config.DB_CONNECTION_STRING
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	dataStore.InitStore(&dataStore.DBStore{DB: db})

	// INITIALIZE SESSION
	sessionStore.InitSessionStore(configVars.Config.SESSION_KEY)

	// ROUTES
	router := mux.NewRouter().StrictSlash(true)

	// SERVE STATIC FILES
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("colorBook/static/"))))

	// RENDER ROUTES
	router.HandleFunc("/login", handlers.RenderLogin).Methods("GET")
	router.HandleFunc("/dashboard", WEB_EnsureAuthenticated(handlers.RenderDashboard)).Methods("GET")
	router.HandleFunc("/logout", WEB_EnsureAuthenticated(handlers.Logout)).Methods("GET")

	// API ENDPOINTS
	router.HandleFunc("/api/get-positions", API_EnsureAuthenticated(api.GetPositions)).Methods("GET")

	router.HandleFunc("/api/register", api.Register).Methods("POST")
	router.HandleFunc("/api/login", api.Login).Methods("POST")

	router.Use(loggingMiddleware)

	// router.HandleFunc("/*")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      c.Handler(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Web server stopped")
}

func main() {
	runServer()
}

/* ----------------- Middleware ------------------ */

func API_EnsureAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")

			// TODO: VERIFY THIS IS THE CASE WITH AXIOS
			if len(bearerToken) == 1 {
				token, error := jwt.Parse(bearerToken[0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(configVars.Config.JWT_TOKEN_SECRET), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(MiddlewareException{ErrorMessage: error.Error(), Success: false})
					return
				}
				if token.Valid {
					// We can also view the token claims here with token.Claims
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(MiddlewareException{ErrorMessage: "Invalid authorization token", Success: false, Redirect: "/login"})
				}
			} else {
				json.NewEncoder(w).Encode(MiddlewareException{ErrorMessage: "Improper authorization header", Success: false, Redirect: "/login"})
			}
		} else {
			json.NewEncoder(w).Encode(MiddlewareException{ErrorMessage: "An authorization header is required", Success: false, Redirect: "/login"})
		}
	})
}

func WEB_EnsureAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userToken, ok := sessionStore.GetSessionValue(w, r, "user_token_string")

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		if len(userToken) > 0 {
			if len(strings.Split(userToken, ".")) == 3 {
				token, error := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(configVars.Config.JWT_TOKEN_SECRET), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(MiddlewareException{ErrorMessage: error.Error(), Success: false})
					return
				}
				if token.Valid {
					next(w, r)
				}
			} else {
				// I would put these redirects in a return but I get
				// a See Other link at the bottom of the page
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL, r.Proto)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
