package cmd

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"os"
)

func prepareDatabase() (*sql.DB, error) {
	db, err := preparePostgres()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getPort() (string, error) {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return utils.EmptyString, errors.New("PORT environment variable missing")
	}
	return port, nil
}

func checkEnv() error {
	_, ok := os.LookupEnv("DBConString")
	if !ok {
		return errors.New("DBConString environment variable missing")
	}
	_, ok = os.LookupEnv("PORT")
	if !ok {
		return errors.New("PORT environment variable missing")
	}
	return nil
}
