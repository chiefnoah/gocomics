package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

//A map would work just as well, but I like things to be hard(ly typed)
type ApiConfig struct {
	UseTLS         bool     `json:"use_tls"`
	ForceTLS       bool     `json:"force_tls"`
	SSLPort        string   `json:"ssl_port"`
	HttpPort       string   `json:"http_port"`
	ComicFolders   []string `json:"comic_folders"`
	DatabaseFolder string   `json:"database_folder"`
}

const CONFIG_FILE string = "config.json"

var globalConfig *ApiConfig

func LoadConfigFile() *ApiConfig {
	if globalConfig != nil {
		return globalConfig
	}
	f, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Println("Could not load config: ", err)
		//TODO: if the file wasn't found, generate default config
	}
	config := ApiConfig{}
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}
	globalConfig = &config

	return globalConfig
}

func WriteConfigFile(config *ApiConfig) error {
	p := filepath.Join(".", CONFIG_FILE)
	json, err := json.MarshalIndent(*config, "", "    ")
	if err != nil {
		log.Println("Unable to marshal settings to JSON: ", err)
		return err
	}

	err = ioutil.WriteFile(p, json, 0644)
	if err != nil {
		log.Println("Unable to write config file: ", err)
		return err
	}
	return nil
}
