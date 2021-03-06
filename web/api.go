package web

import (
	"git.chiefnoah.tech/chiefnoah/gocomics/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

//This is where the REST API stuff will go

//Starts the API server and registers handlers
func Start(c *config.ApiConfig) {
	//TODO: register handlers
	router := gin.Default()

	router.GET("/", rootHandler)
	router.GET("/dbinfo", dbInfoHandler)
	router.GET("/version", versionHandler)
	router.RunTLS(c.HttpPort, "./test.pem", "./test.key")
}

func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "hi")
}

