package main

import (
	"git.chiefnoah.tech/chiefnoah/gocomics/config"
	"git.chiefnoah.tech/chiefnoah/gocomics/scanner"
	"git.chiefnoah.tech/chiefnoah/gocomics/web"

	"git.chiefnoah.tech/chiefnoah/gocomics/database"
	"log"
	"os"
)

//Let's get started!
func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Can't write log file! D:")
	}
	defer f.Close()
	//log.SetOutput(f)
	database.Init()

	c := config.LoadConfigFile()
	go comicscanner.Scan(c.ComicFolders[0])
	//go comicscanner.Scan(c.ComicFolders[1])
	web.Start(c)

}
