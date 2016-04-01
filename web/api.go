package web

import (
	"git.chiefnoah.tech/chiefnoah/gocomics/config"
	"net/http"
	"io"
)

//This is where the REST API stuff will go

//Starts the API server and registers handlers
func Start(c *config.ApiConfig) {
	//TODO: register handlers
	http.HandleFunc("/", RootHandler)

	if c.UseTLS || c.ForceTLS {
		http.ListenAndServeTLS(c.SSLPort, "cert.pem", "key.pem", nil)
	}
	if !c.ForceTLS {
		http.ListenAndServe(c.HttpPort, nil)
	}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "error... nothing here :(")
}