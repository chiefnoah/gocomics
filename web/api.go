package web

import (
	"git.packetlostandfound.us/chiefnoah/gocomics/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//This is where the REST API stuff will go

//Starts the API server and registers handlers
func Start(c *config.ApiConfig) {
	//TODO: register handlers
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/dbinfo", dbInfoHandler)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/comiclist", comicListHandler)
	//router.HandleFunc(`/folders`, foldersHandler)
	//router.HandleFunc(`/folders/`, foldersHandler)
	//router.HandleFunc(`/folders/{root:[0-9]+}`, foldersHandler)
	//router.HandleFunc(`/folders/{root:[0-9]+}/{path:.*}`, foldersHandler)

	log.Printf("Config: %+s", *c)

	if c.UseTLS == true {
		log.Print("Starting HTTPs server...")
		//go func() {
		err := http.ListenAndServeTLS(c.SSLPort, "./test.pem", "./test.key", router)
		if err != nil {
			log.Fatal("Unable to start up HTTPs server: %s", err)
		}
		//}()
	}
	if !c.ForceTLS {
		log.Print("Starting HTTP server...")
		//go func() {
		err := http.ListenAndServe(c.HttpPort, router)
		if err != nil {
			log.Fatal("Unable to start up HTTP server: %s", err)
		}
		//}()
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("Hi"))
}
