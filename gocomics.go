package main

import (
	"git.chiefnoah.tech/chiefnoah/gocomics/config"
	"git.chiefnoah.tech/chiefnoah/gocomics/scanner"
	"git.chiefnoah.tech/chiefnoah/gocomics/web"

	"log"
	"os"
	"strconv"
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
		for i, folder := range c.ComicFolders {
			comicscanner.Scan(folder, strconv.Itoa(i))
		}
	}()

	web.Start(c)
}
