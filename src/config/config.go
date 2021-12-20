package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var Cfg = new(Config)

type Config struct {
	Cookies struct {
		JueJin string `yaml:"juejin"`
	} `yaml:"cookies"`
	Chanify struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"chanify"`
}

func InitConfig(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, Cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
