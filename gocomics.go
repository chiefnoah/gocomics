package main

import (
	"git.chiefnoah.tech/chiefnoah/gocomics/config"
	"git.chiefnoah.tech/chiefnoah/gocomics/web"
	"git.chiefnoah.tech/chiefnoah/gocomics/scanner"

	"git.chiefnoah.tech/chiefnoah/gocomics/database"
)

//Let's get started!
func main() {

	database.Init()

	config := &config.ApiConfig{
		false, false, ":3008", ":3000",
	}
	comicscanner.Scan("./comics")
	web.Start(config)

}
