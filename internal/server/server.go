package server

import (
	"context"
	"net/http"
	"os"

	"tldw-server/internal/event"

	"github.com/gin-gonic/gin"
)

var log = event.Log

// Start the REST API server using the configuration provided
func Start(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// Set http server mode
	gin.SetMode(gin.ReleaseMode)

	// Create router with middleware
	router := gin.Default()

	registerRoutes(router)

	// Get server address from env vars
	downcountAddress := os.Getenv("DOWNCOUNT_ADDRESS")
	downcountPort := os.Getenv("DOWNCOUNT_PORT")
	if downcountAddress == "" || downcountPort == "" {
		panic("server: server address vars not set")
	}
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		log.Infof("http: starting webserver at %s", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("http: web server shutdown complete")
			} else {
				log.Errorf("http: web server closed unexpectedly: %s", err)
			}
		}
	}()

	<-ctx.Done()
	log.Info("http: shutting down web server")

	err := server.Close()
	if err != nil {
		log.Errorf("http: web server shutdown failed: %v", err)
	}
}
