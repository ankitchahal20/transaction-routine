package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pelletier/go-toml"
)

var (
	globalConfig GlobalConfig
)

// Global Configuration
type GlobalConfig struct {
	Database Database `toml:"database"`
	Server   Server   `toml:"server"`
}

// DB configuration
type Database struct {
	Host     string `toml:"host" environment:"POSTGRES_HOST"`
	Port     int    `toml:"port" environment:"POSTGRES_PORT"`
	DBname   string `toml:"dbname" environment:"POSTGRES_DB"`
	User     string `toml:"user" environment:"POSTGRES_USER"`
	Password string `toml:"password" environment:"POSTGRES_PASSWORD"`
}

// server configuration
type Server struct {
	Address      string `toml:"address"`
	ReadTimeOut  int    `toml:"read_time_out"`
	WriteTimeOut int    `toml:"write_time_out"`
}

// Setter method for GlobalConfig
func SetConfig(cfg GlobalConfig) {
	globalConfig = cfg
}

// Getter method for GlobalConfig
func GetConfig() GlobalConfig {
	return globalConfig
}

// Loading the values from default.toml and assigning them as part of GlobalConfig struct
func InitGlobalConfig() error {
	config, err := toml.LoadFile("./config/defaults.toml")
	fmt.Println("Err : ", err)
	if err != nil {
		log.Printf("Error while loading defaults.toml file : %v ", err)
		return err
	}

	var appConfig GlobalConfig
	err = config.Unmarshal(&appConfig)
	if err != nil {
		log.Printf("Error while unmarshalling config : %v", err)
		return err
	}

	// Override with environment variables if they are set
	if host, exists := os.LookupEnv("POSTGRES_HOST"); exists {
		appConfig.Database.Host = host
	}
	if port, exists := os.LookupEnv("POSTGRES_PORT"); exists {
		appConfig.Database.Port, _ = strconv.Atoi(port) // Convert string to int
	}
	if dbname, exists := os.LookupEnv("POSTGRES_DB"); exists {
		appConfig.Database.DBname = dbname
	}
	if user, exists := os.LookupEnv("POSTGRES_USER"); exists {
		appConfig.Database.User = user
	}
	if password, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists {
		appConfig.Database.Password = password
	}

	SetConfig(appConfig)
	return nil
}
