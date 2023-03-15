package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigReader func(filename string) ([]byte, error)
type rawConfig map[string]string

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type config struct {
	serverPort             int
	database               DatabaseConfig
	moviesRowsLimitEnabled bool
	maxMoviesRowsLimit     int64
}

var appConfig config

func Initialize(readFile ConfigReader) {
	configData, err := readFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config file 'config.yml': %v", err)
	}

	configMap := rawConfig{}
	err = yaml.Unmarshal(configData, &configMap)
	if err != nil {
		log.Fatalf("Error reading config file 'config.yml': %v", err)
	}

	appConfig = config{
		serverPort:             configMap.getIntOrPanic("SERVER_PORT"),
		database:               configMap.getDbConfig(),
	}
}

func ServerPort() int {
	return appConfig.serverPort
}

func DbConfig() DatabaseConfig {
	return appConfig.database
}

func (c rawConfig) getStringOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		if _, ok := c[key]; !ok {
			panic(fmt.Sprintf("Config value for '%s' is not set", key))
		}
		value = c[key]
	}
	return value
}

func (c rawConfig) getIntOrPanic(key string) int {
	str := c.getStringOrPanic(key)
	value, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("Config value for '%s' is not a valid int: '%s'", key, str))
	}
	return value
}

func (c rawConfig) getInt64OrPanic(key string) int64 {
	str := c.getStringOrPanic(key)
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Config value for '%s' is not a valid int64: '%s'", key, str))
	}
	return value
}

func (c rawConfig) getBoolOrPanic(key string) bool {
	str := strings.ToLower(c.getStringOrPanic(key))
	if str != "true" && str != "false" {
		panic(fmt.Sprintf("Config value for '%s' is not a valid bool: '%s'", key, str))
	}
	return str == "true"
}

func (c rawConfig) getDbConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     c.getStringOrPanic("DB_HOST"),
		Port:     c.getIntOrPanic("DB_PORT"),
		User:     c.getStringOrPanic("DB_USER"),
		Password: c.getStringOrPanic("DB_PASSWORD"),
		Name:     c.getStringOrPanic("DB_NAME"),
	}
}
