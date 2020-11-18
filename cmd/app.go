package cmd

import (
	"database/sql"
	"ecommerce/router"
	"log"
	"net/http"
)

//App struct
type App struct {
	router.Router
}

//NewApp returns new app struct
func NewApp(db *sql.DB) *App {
	return &App{
		Router: router.NewRouter(db),
	}
}

//Serve to serve the server
func (a *App) Serve() {
	port, err := getPort()
	if err != nil {
		log.Println("Error : Can't find the server address")
		panic(err)
	}
	r := a.Router.Setup()
	log.Println("App : Server is listening")
	http.ListenAndServe("localhost:"+port, r)
}
