package cmd

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"os"

	//Postgres library
	_ "github.com/lib/pq"
)

//DB struct
type DB struct {
	*sql.DB
}

func preparePostgres() (*sql.DB, error) {
	DBUrl, err := getDBUrl()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", DBUrl)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getDBUrl() (string, error) {
	DBUrl, ok := os.LookupEnv("DBConString")
	if !ok {
		return utils.EmptyString, errors.New("DBConString environment variable required")
	}
	return DBUrl, nil
}
