package main

import (
	"git.packetlostandfound.us/chiefnoah/gocomics/config"
	"git.packetlostandfound.us/chiefnoah/gocomics/scanner"
	"git.packetlostandfound.us/chiefnoah/gocomics/web"

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

	c := config.LoadConfigFile()
	go func() {
		for _, folder := range c.ComicFolders {
			comicscanner.Scan(folder)
		}
	}()

	web.Start(c)
}
