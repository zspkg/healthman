package healthman

import (
	"database/sql"
)

const postgresServiceName = "db-postgres"

type dbHealthier struct {
	DB *sql.DB
}

// NewDBHealthier is a constructor of database healthier that implements CheckHealth method based
// on provided database connection
func NewDBHealthier(db *sql.DB) Healthable {
	return dbHealthier{DB: db}
}

// Name returns a determined postgresServiceName value
func (h dbHealthier) Name() string {
	return postgresServiceName
}

// CheckHealth pings a database and returns whether no error occurred
func (h dbHealthier) CheckHealth() bool {
	return h.DB.Ping() == nil
}
