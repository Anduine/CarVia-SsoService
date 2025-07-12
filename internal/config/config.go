package config

import (
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBConnector string
	TokenTTL	time.Duration 	`yaml:"token_ttl" env-required:"true"`
	Port        string   		`yaml:"port"`
	Timeout     time.Duration 	`yaml:"timeout"`
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist: " + path)
	}

	var cfg Config
	data, _ := os.ReadFile(path)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to unmarshal yaml")
	}

	dbConn := getDBconfig()

	cfg.DBConnector = dbConn

	return &cfg
}

func fetchConfigPath() string {
	var result string

	// --config="path/to/config.yaml"
	flag.StringVar(&result, "config", "./configs/config.yaml", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func getDBconfig() string {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		panic("DB_CONN is not set")
	}

	return connStr
}