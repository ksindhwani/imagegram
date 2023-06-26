package app

import (
	"database/sql"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/filesystem"
)

// Dependencies holds the primitives and structs/interfaces that are required
// for the application's business logic.
type Dependencies struct {
	Revision        string
	Config          *config.Config
	DB              *sql.DB
	LocalFileSystem filesystem.FileSystem
}
