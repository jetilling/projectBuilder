package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jetilling/projectBuilder/apiBuilder"
	"github.com/jetilling/projectBuilder/appBuilder"
	"github.com/jetilling/projectBuilder/configVars"
	"github.com/rs/cors"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":8001", "server listen address")
	flag.Parse()

	// INITIALIZE CONFIG VARIABLES
	configVars.InitConfigVars()

	// LOGGING
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Launch appBuilder server is starting...")

	// ROUTES
	router := mux.NewRouter().StrictSlash(true)

	// API ENDPOINTS
	router.HandleFunc("/api/build-project", appBuilder.Build).Methods("POST")
	router.HandleFunc("/api/build-api", apiBuilder.Build).Methods("POST")

	router.Use(loggingMiddleware)

	// router.HandleFunc("/*")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
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
	logger.Println("Launch appBuilder server stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL, r.Proto)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
