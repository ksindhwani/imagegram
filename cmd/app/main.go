package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ksindhwani/imagegram/pkg/app"
	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database/mysql"
	"github.com/ksindhwani/imagegram/pkg/filesystem"
	"github.com/ksindhwani/imagegram/pkg/router"
	"go.uber.org/zap"
)

var (
	revision       = "unknown"
	buildTimestamp = "unknown"
)

func main() {

	// load configuration
	cfg, err := config.New()
	fatalOnError(err, "error loading configuration")

	// initialize persistent stores
	db, err := initializeDB(cfg)
	fatalOnError(err, "error initializing database")

	localFileSystem, err := filesystem.New(filesystem.LOCAL, cfg)
	fatalOnError(err, "error initializing database")

	// initialize application and handlers
	deps := &app.Dependencies{
		Revision:        revision,
		Config:          cfg,
		DB:              db,
		LocalFileSystem: localFileSystem,
	}
	r, err := router.New(deps)
	fatalOnError(err, "could not instantiate router")

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	go func(server *http.Server) {
		log.Print("server running on: %s", cfg.Addr)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen error: %s", err)
		}
	}(server)

	stopCh := make(chan os.Signal, 2)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	log.Print("gracefully shutting down server")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting server down gracefully: %v", err)
	}

}

func fatalOnError(err error, msg string) {
	if err != nil {
		zap.S().Fatalf("%s:%s", msg, err)
	}
}

func initializeDB(cfg *config.Config) (*sql.DB, error) {
	return mysql.NewDB(mysql.ConnectionParams{
		UserID:             cfg.DBUserID,
		Password:           cfg.DBPassword,
		HostName:           cfg.DBHostName,
		Port:               cfg.DBPort,
		Database:           cfg.DBDatabaseName,
		MaxIdleConnections: cfg.DBMaxIdleConnections,
		MaxOpenConnections: cfg.DBMaxOpenConnections,
		MaxConnLifetime:    cfg.DBMaxConnLifetime,
	})
}
