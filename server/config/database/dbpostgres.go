package database

import (
	"database/sql"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	serviceURI := os.Getenv("DB_SERVICE_URI")

	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = os.Getenv("DB_CONNECTION_URI")

	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}
