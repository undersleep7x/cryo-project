package main

//start backend service
import(
	"fmt"
	"log"
	"net/http"

	"github.com/undersleep7x/cryptowallet-v0.1/routes"
	"github.com/undersleep7x/cryptowallet-v0.1/app"
)

func startServer() *http.Server {
	app.InitApp() //kicks off initialization of necessary precursors like redis and logging

	routes.SetupRoutes(app.Router) //register the routes associated with this application

	port := app.Config.App.Port
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: app.Router,
	}
	log.Printf("Server is now running on port %s", port)
	return server

}


func main() {
	server := startServer()
	log.Fatal(server.ListenAndServe())
}