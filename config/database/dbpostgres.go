package database

import (
	"database/sql"
	"fmt"
	//"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 7007
	user     = "cana_user"
	password = "dev123456"
	dbname   = "cana"
)

func Conectar() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, err
}
