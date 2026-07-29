package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTP_Server struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Config struct {
	Env         string      `yaml:"env"`
	StoragePath string      `yaml:"storage_path"`
	HTTPServer  HTTP_Server `yaml:"http_server"`
}

func MustLoad() *Config{
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPathFlag := flag.String("config", "", "path to configuration file")
		flag.Parse()

		configPath = *configPathFlag
	}

	if configPath == "" {
		log.Fatal("no config path set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	log.Println("config loaded successfully", cfg)
	return &cfg
}
