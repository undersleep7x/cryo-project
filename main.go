package main

//start backend service
import (
	"fmt"
	"log"
	"net/http"

	"github.com/undersleep7x/cryo-project/internal/app"
)

func startServer() *http.Server {


	a := app.InitApp() //kicks off initialization of necessary precursors like redis and logging

	port := a.Config.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: a.Router,
	}
	log.Printf("Server is now running on port %s", port)
	return server

}

func main() {
	server := startServer()
	log.Fatal(server.ListenAndServe())
}
