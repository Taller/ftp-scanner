package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database struct {
		Port     string `yaml:"port", envconfig:"DB_PORT"`
		Host     string `yaml:"host", envconfig:"DB_HOST"`
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
		DBName   string `yaml:"dbname", envconfig:"DB_NAME"`
		Schema   string `yaml:"schema", envconfig:"SCHEMA"`
	} `yaml:"database"`
	Ftp struct {
		Host     string `yaml:"host", envconfig:"FTP_HOST"`
		Username string `yaml:"user", envconfig:"FTP_USERNAME"`
		Password string `yaml:"pass", envconfig:"FTP_PASSWORD"`
	} `yaml:"ftp"`
}

func Read(file string) Config {
	var cfg Config
	readFile(file, &cfg)
	readEnv(&cfg)

	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(file string, cfg *Config) {
	f, err := os.Open(file)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
