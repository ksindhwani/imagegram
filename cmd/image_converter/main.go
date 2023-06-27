package main

import (
	"database/sql"
	"log"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/database"
	"github.com/ksindhwani/imagegram/pkg/database/mysql"
	"github.com/ksindhwani/imagegram/pkg/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	fatalOnError(err, "error loading configuration")

	// initialize persistent stores
	db, err := initializeDB(cfg)
	fatalOnError(err, "error initializing database")

	database := database.New(db)
	imageConverterService := service.NewImageConvertorService(cfg, database)
	imageService := service.NewImageService(cfg, database)

	successfulConversions, failedConversions, err := imageConverterService.ConvertImages()
	if err != nil {
		// In Production instead of logging we can log it on log stream
		log.Fatalf("Unable to process images - %w", &err)
	}
	if len(failedConversions) > 0 {
		log.Print("Unable to convert some images")
		for _, image := range failedConversions {
			// In Production instead of logging we can log it on log stream or generate alert
			println(image.ToString())
		}
	}
	err = imageService.UpdateConvertedLocationsForImages(successfulConversions)
	log.Print("Conversion completed")
	fatalOnError(err, "error saving converted image in database")
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
