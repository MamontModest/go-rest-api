package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/mamontmodest/go-rest-api/internal/auth"
	"github.com/mamontmodest/go-rest-api/internal/recipe"
	database "github.com/mamontmodest/go-rest-api/pkg/db"
	"github.com/mamontmodest/go-rest-api/pkg/log"
	"net/http"
	"os"
	"time"
)

var Version = "1.0.0"

var dsn = "postgres://localhost/go_restful?sslmode=disable&user=postgres&password=QWertas1122"

func main() {
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	logger := log.New().With(nil, "version", Version)

	// load application configurations

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	address := fmt.Sprintf(":%v", "8000")
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(&logger, database.NewSDatabase(dsn, "postgres")),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServeTLS(*certFile, *keyFile); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}
func buildHandler(logger *log.Logger, db *database.SDatabase) http.Handler {
	router := routing.New()

	router.Use(
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	recipe.RegisterHandlers(router.Group(""),
		recipe.NewService(recipe.NewRepository(db, logger), logger), logger,
	)
	auth.RegisterHandlers(router.Group(""),
		auth.NewService(auth.NewRepository(db)),
	)

	return router
}
