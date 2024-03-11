package databases

import (
	"database/sql"
	"log/slog"
)

var DbPool *sql.DB

func InitDatabase() error {
	dbtype := "postgresql"
	dbuser := "preguntamebackenduser"
	dbpass := "h4k3rla73n3s4d3n7r0"
	dbhost := "localhost:5432"
	dbname := "preguntame"
	connstr := dbtype + "://" + dbuser + ":" + dbpass + "@" + dbhost + "/" + dbname + "?sslmode=disable";
	slog.Info("Connection string", "connstr", connstr)

	dbPool, err := sql.Open("postgres", connstr)
	if err != nil {
		return err
	}

	err = dbPool.Ping()
	if err != nil {
		return err
	}

	DbPool = dbPool

	return err
}