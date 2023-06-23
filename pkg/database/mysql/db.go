package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// MySQLParseDateLayout is the constant format this library uses
	// as the date layout.
	MySQLParseDateLayout = `%Y-%m-%dT%H:%i:%sZ`
)

// ConnectionParams is the configuration that would be need for
// every single Database *tx.DB initialization.
type ConnectionParams struct {
	UserID             string
	Password           string
	HostName           string
	Port               int
	Database           string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxConnLifetime    time.Duration
}

// NewDB returns a new instance of mysql DB which is a *sql.DB
// wrapped under *tx.DB which provides modular functional
// transactional capabilities.
func NewDB(cfg ConnectionParams) (*sql.DB, error) {
	connectionString := cfg.ResolveURL()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

	return db, nil
}

// ResolveURL takes in ConnectionParams to provide a connection string URL
// that looks like so: "root:root@tcp(127.0.0.1:3306)/canvas?parseTime=true".
func (c ConnectionParams) ResolveURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.UserID,
		c.Password,
		c.HostName,
		c.Port,
		c.Database)
}
